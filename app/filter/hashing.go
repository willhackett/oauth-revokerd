package filter

import (
	"github.com/spaolacci/murmur3"
)

// generateHashes returns 4 64-bit - 2 x 128-bit MurmurHash3 hashes
func generateHashes(data []byte) [4]uint64 {
	h1, h2 := murmur3.Sum128(data)
	buff := make([]byte, len(data)+1)
	copy(buff, data)
	buff[len(data)] = byte(1)
	h3, h4 := murmur3.Sum128(buff)
	return [4]uint64{h1, h2, h3, h4}
}

// getHash retrieves the n round of hashing from 4 pregenerated hashes
func getHash(hash [4]uint64, n uint64) uint64 {
	index := 2 + (((n + (n % 2)) % 4) / 2)
	pre := hash[n%2]
	post := n * hash[index]
	return pre + post
}
