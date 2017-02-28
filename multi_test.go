package bloom

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax
	// characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

func BenchmarkBloom1(b *testing.B) {

	var (
		ids  []string
		locs = make([]uint64, 100)
	)
	const (
		e = 20000
	)
	// create ids
	for i := 0; i < e*2; i++ {
		ids = append(ids, RandStringBytesMaskImprSrc(20))
	}
	type TestCase struct {
		k     int
		n     int
		m     int
		f     int
		multi bool
	}

	tests := []TestCase{
		{5, 1024, 1024 * 8 * 2, 20, false},
		{6, 1024, 1024 * 8, 20, false},
		{10, 1024, 1024 * 8, 20, false},
		{5, 516, 1024 * 8, 20, false},
		{5, 1024, 1024 * 8, 20, false},
		{5, 1024, 1024 * 8, 15, false},
		{5, 1024, 1024 * 8, 10, false},
		{5, 1024, 1024 * 8, 5, false},
		{5, 1024, 1024 * 8, 1, false},
		{5, 1024, 1024 * 8, 20, true},
		{5, 1024, 1024 * 8, 15, true},
		{5, 1024, 1024 * 8, 10, true},
		{5, 1024, 1024 * 8, 5, true},
		{5, 1024, 1024 * 8, 1, true},
		{5, 1024, 1024 * 10 * 8, 20, true},
	}

	for _, tc := range tests {
		b.Run(fmt.Sprintf("%v", tc), func(b *testing.B) {

			// prepare
			var bfs []*BloomFilter
			// create 16 bloom filters
			var j int
			for i := 0; i < tc.f; i++ {
				bf := New(uint(tc.m), uint(tc.k))
				bfs = append(bfs, bf)

				for k := 0; k < tc.n; k++ {
					s := ids[j%len(ids)]
					j++
					bf.Add([]byte(s))
				}
			}

			b.ResetTimer()
			for l := 0; l < b.N; l++ {
				contains := 0
				for i := 0; i < len(ids); i++ {
					s := ids[i%len(ids)]
					if tc.multi {
						if MultiTest([]byte(s), locs, bfs) {
							contains++
						}
					} else {
						for _, bf := range bfs {
							if bf.Test([]byte(s)) {
								contains++
							}
						}
					}
				}
				//			fmt.Println("contains:", contains, "should contain:", tc.n*tc.f)
			}
		})
	}
}

func TestMultiTest(t *testing.T) {
	var (
		m    uint = 1024
		k    uint = 5
		v1        = []byte("value")
		v2        = []byte("value2")
		v3        = []byte("value3")
		locs      = make([]uint64, 100)
	)
	bf := New(m, k)
	bf.Add(v1)
	if !bf.Test(v1) {
		t.Fail()
	}
	if bf.Test(v2) {
		t.Fail()
	}

	bf2 := New(m, k)
	bf2.Add(v2)
	if !bf2.Test(v2) {
		t.Fail()
	}
	if bf2.Test(v1) {
		t.Fail()
	}

	if !MultiTest(v1, locs, []*BloomFilter{bf, bf2}) {
		t.Fail()
	}

	if !MultiTest(v2, locs, []*BloomFilter{bf, bf2}) {
		t.Fail()
	}
	if MultiTest(v3, locs, []*BloomFilter{bf, bf2}) {
		t.Fail()
	}
}
