package rando

import (
	"math/rand"
	"sync"
)

var pool = &sipPool{p: sync.Pool{New: func() interface{} {
	return NewSource()
}}}

// Pooled returns global source used in implementation of global methods of this package.
// Global source uses sync.Pool, so it is quite cheap. But if you need to generate
// many values (>10) at once, create new source with New().
func Pooled() rand.Source {
	return pool
}

type sipPool struct {
	p sync.Pool
}

func (s *sipPool) Uint64() uint64 {
	p := s.p.Get().(*sipRand)
	v := p.Uint64()
	s.p.Put(p)
	return v
}

func (c *sipPool) Int63() int64 {
	return int64(c.Uint64() >> 1)
}

func (c *sipPool) Seed(i int64) {}

var globalRand = rand.New(pool)

// Int63 returns a non-negative random 63-bit integer as an int64
// from the default Source.
func Int63() int64 { return pool.Int63() }

// Uint32 returns a random 32-bit value as a uint32
// from the default Source.
func Uint32() uint32 { return globalRand.Uint32() }

// Uint64 returns a random 64-bit value as a uint64
// from the default Source.
func Uint64() uint64 { return pool.Uint64() }

// Int31 returns a non-negative random 31-bit integer as an int32
// from the default Source.
func Int31() int32 { return globalRand.Int31() }

// Int returns a non-negative random int from the default Source.
func Int() int { return globalRand.Int() }

// Int63n returns, as an int64, a non-negative random number in [0,n)
// from the default Source.
// It panics if n <= 0.
func Int63n(n int64) int64 { return globalRand.Int63n(n) }

// Int31n returns, as an int32, a non-negative random number in [0,n)
// from the default Source.
// It panics if n <= 0.
func Int31n(n int32) int32 { return globalRand.Int31n(n) }

// Intn returns, as an int, a non-negative random number in [0,n)
// from the default Source.
// It panics if n <= 0.
func Intn(n int) int { return globalRand.Intn(n) }

// Float64 returns, as a float64, a random number in [0.0,1.0)
// from the default Source.
func Float64() float64 { return globalRand.Float64() }

// Float32 returns, as a float32, a random number in [0.0,1.0)
// from the default Source.
func Float32() float32 { return globalRand.Float32() }

// Perm returns, as a slice of n ints, a random permutation of the integers [0,n)
// from the default Source.
func Perm(n int) []int { return globalRand.Perm(n) }

// Shuffle randomizes the order of elements using the default Source.
// n is the number of elements. Shuffle panics if n < 0.
// swap swaps the elements with indexes i and j.
func Shuffle(n int, swap func(i, j int)) { globalRand.Shuffle(n, swap) }

// Read generates len(p) random bytes from the default Source and
// writes them into p. It always returns len(p) and a nil error.
// Read, unlike the Rand.Read method, is safe for concurrent use.
func Read(p []byte) (n int, err error) { return globalRand.Read(p) }

// NormFloat64 returns a normally distributed float64 in the range
// [-math.MaxFloat64, +math.MaxFloat64] with
// standard normal distribution (mean = 0, stddev = 1)
// from the default Source.
// To produce a different normal distribution, callers can
// adjust the output using:
//
//  sample = NormFloat64() * desiredStdDev + desiredMean
//
func NormFloat64() float64 { return globalRand.NormFloat64() }

// ExpFloat64 returns an exponentially distributed float64 in the range
// (0, +math.MaxFloat64] with an exponential distribution whose rate parameter
// (lambda) is 1 and whose mean is 1/lambda (1) from the default Source.
// To produce a distribution with a different rate parameter,
// callers can adjust the output using:
//
//  sample = ExpFloat64() / desiredRateParameter
//
func ExpFloat64() float64 { return globalRand.ExpFloat64() }
