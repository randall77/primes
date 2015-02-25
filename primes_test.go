package primes

import "testing"

func isPrime(p int) bool {
	if p < 2 {
		return false
	}
	for d := 2; d < p; d++ {
		if p%d == 0 {
			return false
		}
	}
	return true
}

func TestSmallPrimes(t *testing.T) {
	s := &Set{}
	for p := -10; p < 10000; p++ {
		if isPrime(p) != s.Contains(p) {
			t.Errorf("prime(%d) want %v, got %v", p, isPrime(p), s.Contains(p))
		}
	}
}

// pi(n) = # of primes <= n
var pi = map[int]int{
	10:         4,
	100:        25,
	1000:       168,
	10000:      1229,
	100000:     9592,
	1000000:    78498,
	10000000:   664579,
	100000000:  5761455,
	1000000000: 50847534,
}

func TestPrimeCounts(t *testing.T) {
	testPrimeCounts(t, &Set{})
}
func TestPrimeCounts2(t *testing.T) {
	testPrimeCounts(t, New(1e9))
}
func testPrimeCounts(t *testing.T, s *Set) {
	for n, c := range pi {
		d := 0
		for i := 0; i < n; i++ {
			if s.Contains(i) {
				d++
			}
		}

		if c != d {
			t.Errorf("pi(%d)=%d, got %d\n", n, c, d)
		}
	}
}

func BenchmarkPrime1e6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New(1e6).Contains(1e6)
	}
}

func BenchmarkPrime1e9(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New(1e9).Contains(1e9)
	}
}

func BenchmarkContains(b *testing.B) {
	s := &Set{}
	for i := 0; i < b.N; i++ {
		s.Contains(1e6)
	}
}
