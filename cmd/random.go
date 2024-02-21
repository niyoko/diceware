package cmd

import (
	"crypto/rand"
	"math/bits"
)

type rndBit struct {
	generated []byte
	used      int
}

func (r *rndBit) next(n int) (result uint64) {
	if n > 64 {
		panic("n > 64")
	}

	if r.used+n > len(r.generated)*8 {
		r.fillBytes()
	}

	for i := 0; i < n; i++ {
		byteIndex := r.used / 8
		bitIndex := r.used % 8
		if r.generated[byteIndex]&(1<<bitIndex) != 0 {
			result |= 1 << i
		}
		r.used++
	}

	return
}

func (r *rndBit) randN(n int) (result uint64) {
	if n <= 0 {
		panic("n <= 0")
	}

	if n == 1 {
		return 0
	}

	iterN := 0
	bit := bits.Len(uint(n - 1))
	for {
		if iterN > 5000 {
			panic("iterN > 5000")
		}

		rnd := r.next(bit)
		if rnd < uint64(n) {
			return rnd
		}

		iterN++
	}
}

func (r *rndBit) fillBytes() {
	gen := make([]byte, 24)
	_, err := rand.Read(gen)
	if err != nil {
		panic(err)
	}

	r.generated = append(r.generated, gen...)
}
