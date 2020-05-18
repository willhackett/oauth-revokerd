package filter

import (
	"encoding/json"
	"errors"
	"math"
	"sync"
)

var (
	errElements      = errors.New("error: elements must be greater than 0")
	errFalsePositive = errors.New("error: falsePositive must be greater than 0 and less than 1")
)

type ExportInterface struct {
	size uint64  `json:"size"`
	bits []uint8 `json:"bits"`
	hash uint64  `json:"hash"`
}

// BloomFilter contains the information for a BloomFilter data store.
type BloomFilter struct {
	size  uint64        // number of bits (bit array is size/8+1)
	bits  []uint8       // main bit array
	hash  uint64        // number of hash rounds
	mutex *sync.RWMutex // mutex for locking Add, Test, and Reset operations
}

// Add adds the data to the BloomFilter.
func (bf *BloomFilter) Add(data []byte) {
	// generate hashes
	hashes := generateHashes(data)
	bf.mutex.Lock()
	for i := uint64(0); i < bf.hash; i++ {
		index := getHash(hashes, i) % bf.size
		bf.bits[index/8] |= (1 << (index % 8))
	}
	bf.mutex.Unlock()
}

// Reset clears the BloomFilter.
func (bf *BloomFilter) Reset() {
	bf.mutex.Lock()
	bf.bits = make([]uint8, bf.size/8+1)
	bf.mutex.Unlock()
}

// Test returns a bool if the data is in the BloomFilter. True indicates that the data
// may be in the BloomFilter, while false indicates that the data is not in the BloomFilter.
func (bf *BloomFilter) Test(data []byte) bool {
	hashes := generateHashes(data)
	bf.mutex.RLock()
	defer bf.mutex.RUnlock()
	for i := uint64(0); i < uint64(bf.hash); i++ {
		index := getHash(hashes, i) % bf.size
		// check if index%8-th bit is not active
		if (bf.bits[index/8] & (1 << (index % 8))) == 0 {
			return false
		}
	}
	return true
}

// Export the filter as JSON
func (bf *BloomFilter) Export() ([]byte, error) {
	exp := &ExportInterface{
		bf.size,
		bf.bits,
		bf.hash,
	}

	d, err := json.Marshal(exp)
	if err != nil {
		return []byte(""), err
	}

	return d, nil
}

// New creates a bloom filter
func New(size int, threshold float64) (*BloomFilter, error) {
	bf := new(BloomFilter)
	// number of bits
	m := (-1 * float64(size) * math.Log(threshold)) / math.Pow(math.Log(2), 2)
	// number of hash operations
	k := (m / float64(size)) * math.Log(2)

	bf.mutex = &sync.RWMutex{}
	bf.size = uint64(math.Ceil(m))
	bf.hash = uint64(math.Ceil(k))
	bf.bits = make([]uint8, bf.size/8+1)
	return bf, nil
}
