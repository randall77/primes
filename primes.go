/*
Package primes is used for computing lots of primes quickly.  It can compute
primes up to 1e9 in less than a second.  It uses a byte of memory for every 30
numbers it checks, so it can compute primes up to 30e9 with 1GB of memory.

Allocate a prime Set using s := &Set{} or s := New(maxprime) and query it with s.Contains(1001).

A Set instance is not thread safe.

Originally developed to solve https://projecteuler.net/problem=501
*/

package primes

// A slot encodes one of the numbers 1, 7, 11, 13, 17, 19, 23, and 29 (mod 30) using values 0-7.
// Slot numbers are those that are nonzero mod 2, 3, and 5.
type slot byte

// look up slot values.
var slot2int = [8]byte{1, 7, 11, 13, 17, 19, 23, 29}

// mask30[x] = bit mask encoding slot, or 0 if index is not a slot
var mask30 = [30]byte{0, 1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 4, 0, 8, 0, 0, 0, 16, 0, 32, 0, 0, 0, 64, 0, 0, 0, 0, 0, 128}

// A Set is a lazily computed set of all prime numbers.
type Set struct {
	// bit j of sieve[i] encodes the primeness of i*30+slot2int[j]
	sieve []byte
}

// New allocates a new set.  n is the expected maximum number you'll ever ask about.
func New(n int) *Set {
	return &Set{sieve: make([]byte, 0, (n+chunk)/30)}
}

// Contains returns whether p is prime.
func (s *Set) Contains(p int) bool {
	if p < 0 {
		return false
	}
	if p == 2 || p == 3 || p == 5 {
		return true
	}
	b := p / 30
	for b >= len(s.sieve) {
		s.computeMore()
	}
	return s.sieve[b]&mask30[p-b*30] != 0
}

// For sieving, we are given i and we need to mark all values i+n*p as
// composite, n >= 0.  We break n up into j+30*k, 0<=j<30, so that
// iterating over k marks the same slot in every byte.  Only 8 of the
// 30 j's are worth doing (the others map to values which are not
// slots).  This table finds those j's and the corresponding masks we
// need.
type iterEntry struct {
	j    byte
	mask byte
}

var iterInfo [30][8][8]iterEntry

func init() {
	for i := 0; i < 30; i++ {
		for p := 0; p < 8; p++ {
			n := 0
			for j := 0; j < 30; j++ {
				x := i + j*int(slot2int[p])
				if x%2 != 0 && x%3 != 0 && x%5 != 0 {
					iterInfo[i][p][n] = iterEntry{byte(j), ^mask30[x%30]}
					n++
				}
			}
			if n != 8 {
				panic("bad j count")
			}
		}
	}
}

// A chunk is the range of values that are sieved for primes all at
// once.  A chunk occupies chunk/30 bytes.  We want this size to fit
// in cache because we access it in a strided fashion.
const chunk = 30 * (1 << 15) // 32KB of data

// A chunk with all bits set
var one [chunk / 30]byte

func init() {
	for i := range one {
		one[i] = 0xff
	}
}

func (set *Set) computeMore() {
	// extend sieve array by one chunk
	start := len(set.sieve) * 30
	end := start + chunk
	set.sieve = append(set.sieve, one[:]...)

	if start == 0 {
		// For bootstrapping, initialize first byte.  1 is not prime, the rest are.
		set.sieve[0] = 0xfe
	}

	// Iterate over primes we know so far, mark their multiples in the new chunk as not prime.
	sieve := set.sieve
	for base, mask := range sieve {
		for s := slot(0); mask != 0; s++ {
			if mask&1 == 0 { // TODO: find first bit instruction?
				mask >>= 1
				continue
			}
			mask >>= 1

			// The next prime to sieve with.
			p := base*30 + int(slot2int[s])

			// Mark all multiples of p, starting at p*p, as composite.
			i := p * p
			if i >= end {
				return
			}
			if i < start {
				// set i to the first multiple of p which is >= start
				i = (start + p - 1) / p * p
			}

			row := &iterInfo[i%30][s]
			var indexes [8]int
			var masks [8]byte
			for k := 0; k < 8; k++ {
				indexes[k] = (i + int(row[k].j)*p) / 30
				masks[k] = row[k].mask
			}

			inner(sieve, &indexes, masks, p)
		}
	}
}
