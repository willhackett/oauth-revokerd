package filter

import (
	"github.com/willf/bloom"
)

type Filter struct {
	bf *bloom.BloomFilter
}

// New returns a Filter with a _len_ and 5 base hashes
func New(len uint) *Filter {
	bf := bloom.New(len, 5)
	return &Filter{
		bf,
	}
}

// Add appends a string to the filter
func (f *Filter) Add(jti string) {
	f.bf = f.bf.AddString(jti)
}

// Test returns true if the string is in the BloomFilter, false otherwise.
// If true, the result might be a false positive. If false, the data
// is definitely not in the set.
func (f *Filter) Test(data string) bool {
	return f.bf.Test([]byte(data))
}

// MarshalJSON implements json.Marshaler interface.
func (f *Filter) MarshalJSON() ([]byte, error) {
	return f.bf.MarshalJSON()
}
