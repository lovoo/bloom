package bloom

// test based on locations
func (f *BloomFilter) test(locs []uint64) bool {
	for i := uint(0); i < f.k; i++ {
		if !f.b.Test(uint(locs[i] % uint64(f.m))) {
			return false
		}
	}
	return true
}

// location returns the ith hashed location using the four base hash values
func location(h [4]uint64, i uint) uint64 {
	ii := uint64(i)
	return h[ii%2] + ii*h[2+(((ii+(ii%2))%4)/2)]
}

// MultiTest returns true if one bloom filter contains the data, false
// otherwise.
func MultiTest(data []byte, locs []uint64, bfs []*BloomFilter) bool {
	// assume same size
	var k uint
	for _, f := range bfs {
		if f.k > k {
			k = f.k
		}
	}

	// calculate locations
	h := baseHashes(data)
	for i := uint(0); i < k; i++ {
		locs[i] = location(h, i)
	}

	// test in each bloom filter
	for _, bf := range bfs {
		if bf.test(locs) {
			return true
		}
	}
	return false
}
