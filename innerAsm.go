// +build amd64

package primes

//go:noescape
func inner(sieve []byte, indexes *[8]int, masks [8]byte, p int)
