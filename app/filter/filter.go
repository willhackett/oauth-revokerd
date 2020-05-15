package filter

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"sync"
)

var (
	errElements      = errors.New("error: elements must be greater than 0")
	errFalsePositive = errors.New("error: falsePositive must be greater than 0 and less than 1")
)

// BloomFilter contains the information for a BloomFilter data store.
type BloomFilter struct {
	size  uint64        // number of bits (bit array is size/8+1)
	bits  []uint8       // main bit array
	hash  uint64        // number of hash rounds
	mutex *sync.RWMutex // mutex for locking Add, Test, and Reset operations
}

// Init initializes and returns a new BloomFilter, or an error. Given a number of
// elements, it accurately states if data is not added. Within a falsePositive
// rate, it will indicate if the data has been added.
func Init(elements int, falsePositive float64) (*BloomFilter, error) {
	if elements <= 0 {
		return nil, errElements
	}
	if falsePositive <= 0 || falsePositive >= 1 {
		return nil, errFalsePositive
	}

	r := BloomFilter{}
	// number of bits
	m := (-1 * float64(elements) * math.Log(falsePositive)) / math.Pow(math.Log(2), 2)
	// number of hash operations
	k := (m / float64(elements)) * math.Log(2)

	r.mutex = &sync.RWMutex{}
	r.size = uint64(math.Ceil(m))
	r.hash = uint64(math.Ceil(k))
	r.bits = make([]uint8, r.size/8+1)
	return &r, nil
}

// Add adds the data to the BloomFilter.
func (r *BloomFilter) Add(data []byte) {
	// generate hashes
	hash := generateMultiHash(data)
	r.mutex.Lock()
	for i := uint64(0); i < r.hash; i++ {
		index := getRound(hash, i) % r.size
		r.bits[index/8] |= (1 << (index % 8))
	}
	r.mutex.Unlock()
}

// Reset clears the BloomFilter.
func (r *BloomFilter) Reset() {
	r.mutex.Lock()
	r.bits = make([]uint8, r.size/8+1)
	r.mutex.Unlock()
}

// Test returns a bool if the data is in the BloomFilter. True indicates that the data
// may be in the BloomFilter, while false indicates that the data is not in the BloomFilter.
func (r *BloomFilter) Test(data []byte) bool {
	// generate hashes
	hash := generateMultiHash(data)
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for i := uint64(0); i < uint64(r.hash); i++ {
		index := getRound(hash, i) % r.size
		// check if index%8-th bit is not active
		if (r.bits[index/8] & (1 << (index % 8))) == 0 {
			return false
		}
	}
	return true
}

// Merges the sent BloomFilter into itself.
func (r *BloomFilter) Merge(m *BloomFilter) error {
	if r.size != m.size || r.hash != m.hash {
		return errors.New("BloomFilters must have the same m/k parameters")
	}

	r.mutex.Lock()
	m.mutex.RLock()
	for i := 0; i < len(m.bits); i++ {
		r.bits[i] |= m.bits[i]
	}
	r.mutex.Unlock()
	m.mutex.RUnlock()
	return nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (r *BloomFilter) MarshalBinary() ([]byte, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	out := make([]byte, len(r.bits)+17)
	// store a version for future compatibility
	out[0] = 1
	binary.BigEndian.PutUint64(out[1:9], r.size)
	binary.BigEndian.PutUint64(out[9:17], r.hash)
	copy(out[17:], r.bits)
	return out, nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (r *BloomFilter) UnmarshalBinary(data []byte) error {
	// 17 bytes for version + size + hash and 1 byte at least for bits
	if len(data) < 17+1 {
		return fmt.Errorf("incorrect length: %d", len(data))
	}
	if data[0] != 1 {
		return fmt.Errorf("unexpected version: %d", data[0])
	}
	if r.mutex == nil {
		r.mutex = new(sync.RWMutex)
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.size = binary.BigEndian.Uint64(data[1:9])
	r.hash = binary.BigEndian.Uint64(data[9:17])
	// sanity check against the bits being the wrong size
	if len(r.bits) != int(r.size/8+1) {
		r.bits = make([]uint8, r.size/8+1)
	}
	copy(r.bits, data[17:])
	return nil
}
