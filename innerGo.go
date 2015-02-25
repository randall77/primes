// +build 386 amd64p32 arm ppc64 ppc64le

package primes

func inner(sieve []byte, indexes *[8]int, masks [8]byte, p int) {
	for k := 0; k < 8; k++ {
		m := masks[k]

		for i := indexes[k]; i < len(sieve); i += p {
			sieve[i] &= m
		}
	}
}
