// rando is a (practically) secure and (practically) fast random generator.
//
// Technically, it is entropy "multiplicator" that takes crypto/rand as source
// of entropy and then applies high quality permutation.
//
// Current implementation:
//
// To produce each 128bit of output generator:
//   - xor state with:
//   -- 64bit per generator counter
//   -- 64bit per generator unique id
//   -- 64bit salt randomly taken from global entropy pool of 128 elements.
//   - performs 5 rounds of SipHash permutaion
//   - takes two elements and xor them with per-generator mask
// Generator is initialized in following way:
//   - id is get by permuting incremented global value.
//     Incrementing global value and using invertable permutation, it quaranteed,
//     that every generator in this process has unique id.
//   - counter is initialized with global value xor-ed with id
//   - small generators, that produces indices to entropy pool, is derived from id in
//     non-linear way.
//   - permutaion state and per-generator mask is then randomly taken from entropy pool.
//
// Entropy pool of 128 uint64 is refreshed 10 times a second from crypto/rand.
// All other global variables are also initialized from crypto/rand at process start.
//
//		import "fmt"
//		import "math/rand"
//		import rando "github.com/funny-falcon/go-rando"
//
//		func main() {
//			// global methods are same to math/rand
//			fmt.Printf("%x %d\n", rando.Uint64(), rando.Intn(1000))
//			// for generating many values it is better to wrap it to rand.New
//			// Note: it is much cheaper to call rando.NewSource() than rand.NewSource()
//			srng := rand.New(rando.New())
//			fmt.Printf("%x %d\n", srng.Uint64(), srng.Intn(1000))
//		}
//
// Generator uses permutation from cryptographic function, and it avoids presumable
// pitfal of that permutation, so it is expected to pass any statictical tests.
// It seems that the expectations are justified: dieharder, practrand and BigCrush
// were completed without any fail.
//
// Despite being practically secure, Rando is almost as fast as math/rand:
//
// - it is much smaller to allocate and initialize, because default math/rand
// source consists of 637 64bit numbers, and they all have to be initialized
// on seeding.
//
// - it is just ~10% slower in long term. I don't know really why, may be still
// because of smaller footprint, and because it does less slice accesses.
//
package rando

import (
	"math/rand"
	"sync/atomic"
)

type sipRand struct {
	// siphash state
	v0, v1, v2, v3 uint64
	// unique generator id
	id uint64
	// per-generator counter
	cnt uint64
	// whitening of output
	w0, w1 uint64
	// generator used to take random items from a pool
	pos uint32
	// signal, should we iterate, or just take second 128bit value
	f bool
}

// Just calls math/rand.New for convenience of drop-in replacement
func New(src rand.Source) *rand.Rand {
	return rand.New(src)
}

// NewSource returns properly initialized rand.Source instance which produce practically
// secure random numbers.
// It is not goroutine safe, but it is quite cheap to create new for every 10 generate values.
// If you need to generate just 1-2 values, use global methods of this package.
func NewSource() rand.Source {
	s := &sipRand{}
	// Initialize counter at some point.
	// It could be safely initialized as 1, but let it be different at least per-process.
	s.cnt = genoffset
	// Take unique generator id, derive seed generator from it.
	id := atomic.AddUint64(&genid, genidadd)
	// Note: "seed generator" opens side channel on id by timing access to seedspace.
	// But even if side channel is used, there is still 32bit of unknown id bits.
	pos := uint32(id)
	// permute id a bit
	id ^= id >> 32
	id = id*genidmult + genidadd
	id ^= id >> 32
	s.id = id
	s.cnt ^= id
	pos += uint32(id)
	s.pos = pos

	m, a := pos>>8&^2|5, pos>>16|1
	seed := func() uint64 {
		pos = pos*m + a
		return atomic.LoadUint64(&seedspace[pos&0x7f])
	}

	// s.v0 stays zero, because it will be xor-ed in first permutation
	s.v1, s.v2, s.v3 = seed(), seed(), seed()
	s.w0, s.w1 = seed(), seed()
	s.permute()
	s.f = true
	return s
}

func (s *sipRand) permute() {
	v0, v1, v2, v3 := s.v0, s.v1, s.v2, s.v3
	v0 ^= s.frompool()
	v1 ^= s.id
	v3 ^= s.cnt
	s.cnt += uint64(s.pos)
	for i := 0; i < 5; i++ {
		v0 += v1
		v1 = v1<<13 | v1>>(64-13)
		v1 ^= v0
		v0 = v0<<32 | v0>>(64-32)

		v2 += v3
		v3 = v3<<16 | v3>>(64-16)
		v3 ^= v2

		v0 += v3
		v3 = v3<<21 | v3>>(64-21)
		v3 ^= v0

		v2 += v1
		v1 = v1<<17 | v1>>(64-17)
		v1 ^= v2
		v2 = v2<<32 | v2>>(64-32)
	}
	s.v0, s.v1, s.v2, s.v3 = v0, v1, v2, v3
}

func (s *sipRand) frompool() uint64 {
	pos := s.pos
	s.pos = pos*9 + 1
	return atomic.LoadUint64(&seedspace[(pos>>25)&0x7f])
}

func (s *sipRand) Uint64() uint64 {
	if s.f {
		s.f = false
		return s.v1 ^ s.w0
	} else {
		s.f = true
		res := s.v3 ^ s.w1
		s.permute()
		return res
	}
}

func (c *sipRand) Int63() int64 {
	return int64(c.Uint64() >> 1)
}

func (c *sipRand) Seed(i int64) {}
