package rando_test

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	rando "github.com/funny-falcon/go-rando"
	"github.com/funny-falcon/go-rando/crando"
)

// almost copy of math lockedSource
type StdGlob struct {
	sync.Mutex
	rand.Source64
}

var stdGlob = StdGlob{Source64: rand.NewSource(time.Now().UnixNano()).(rand.Source64)}

func (s *StdGlob) Uint64() uint64 {
	s.Lock()
	n := s.Source64.Uint64()
	s.Unlock()
	return n
}

func (s *StdGlob) Int63() int64 {
	s.Lock()
	n := s.Source64.Int63()
	s.Unlock()
	return n
}

func (s *StdGlob) Seed(i int64) {
	s.Lock()
	s.Seed(i)
	s.Unlock()
}

func NewStdGlob() rand.Source {
	return &stdGlob
}

func NewStd() rand.Source {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func bench(b *testing.B, source func() rand.Source, n int) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r := rand.New(source())
			for i := 0; i < n; i++ {
				r.Int63n(100)
			}
			b.SetBytes(8 * int64(n))
		}
	})
}

func BenchmarkMathRand_1(b *testing.B)     { bench(b, NewStd, 1) }
func BenchmarkMathRand_10(b *testing.B)    { bench(b, NewStd, 10) }
func BenchmarkMathRand_100(b *testing.B)   { bench(b, NewStd, 100) }
func BenchmarkMathRand_1000(b *testing.B)  { bench(b, NewStd, 1000) }
func BenchmarkMathRand_10000(b *testing.B) { bench(b, NewStd, 10000) }

func BenchmarkRando_1(b *testing.B)     { bench(b, rando.NewSource, 1) }
func BenchmarkRando_10(b *testing.B)    { bench(b, rando.NewSource, 10) }
func BenchmarkRando_100(b *testing.B)   { bench(b, rando.NewSource, 100) }
func BenchmarkRando_1000(b *testing.B)  { bench(b, rando.NewSource, 1000) }
func BenchmarkRando_10000(b *testing.B) { bench(b, rando.NewSource, 10000) }

func BenchmarkCryptoRand_1(b *testing.B)     { bench(b, crando.NewCrand, 1) }
func BenchmarkCryptoRand_10(b *testing.B)    { bench(b, crando.NewCrand, 10) }
func BenchmarkCryptoRand_100(b *testing.B)   { bench(b, crando.NewCrand, 100) }
func BenchmarkCryptoRand_1000(b *testing.B)  { bench(b, crando.NewCrand, 1000) }
func BenchmarkCryptoRand_10000(b *testing.B) { bench(b, crando.NewCrand, 10000) }

func BenchmarkMathRandGlobal_1(b *testing.B)     { bench(b, NewStdGlob, 1) }
func BenchmarkMathRandGlobal_10(b *testing.B)    { bench(b, NewStdGlob, 10) }
func BenchmarkMathRandGlobal_100(b *testing.B)   { bench(b, NewStdGlob, 100) }
func BenchmarkMathRandGlobal_1000(b *testing.B)  { bench(b, NewStdGlob, 1000) }
func BenchmarkMathRandGlobal_10000(b *testing.B) { bench(b, NewStdGlob, 10000) }

func BenchmarkRandoGlobal_1(b *testing.B)     { bench(b, rando.Pooled, 1) }
func BenchmarkRandoGlobal_10(b *testing.B)    { bench(b, rando.Pooled, 10) }
func BenchmarkRandoGlobal_100(b *testing.B)   { bench(b, rando.Pooled, 100) }
func BenchmarkRandoGlobal_1000(b *testing.B)  { bench(b, rando.Pooled, 1000) }
func BenchmarkRandoGlobal_10000(b *testing.B) { bench(b, rando.Pooled, 10000) }

func BenchmarkCryptoGlobal_1(b *testing.B)     { bench(b, crando.Pooled, 1) }
func BenchmarkCryptoGlobal_10(b *testing.B)    { bench(b, crando.Pooled, 10) }
func BenchmarkCryptoGlobal_100(b *testing.B)   { bench(b, crando.Pooled, 100) }
func BenchmarkCryptoGlobal_1000(b *testing.B)  { bench(b, crando.Pooled, 1000) }
func BenchmarkCryptoGlobal_10000(b *testing.B) { bench(b, crando.Pooled, 10000) }

// Test just math.Source, without math.Rand

func benchsrc(b *testing.B, source func() rand.Source, n int) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r := source().(rand.Source64)
			for i := 0; i < n; i++ {
				r.Uint64()
			}
			b.SetBytes(8 * int64(n))
		}
	})
}

func BenchmarkMathRandSrc_1(b *testing.B)     { benchsrc(b, NewStd, 1) }
func BenchmarkMathRandSrc_10(b *testing.B)    { benchsrc(b, NewStd, 10) }
func BenchmarkMathRandSrc_100(b *testing.B)   { benchsrc(b, NewStd, 100) }
func BenchmarkMathRandSrc_1000(b *testing.B)  { benchsrc(b, NewStd, 1000) }
func BenchmarkMathRandSrc_10000(b *testing.B) { benchsrc(b, NewStd, 10000) }

func BenchmarkRandoSrc_1(b *testing.B)     { benchsrc(b, rando.NewSource, 1) }
func BenchmarkRandoSrc_10(b *testing.B)    { benchsrc(b, rando.NewSource, 10) }
func BenchmarkRandoSrc_100(b *testing.B)   { benchsrc(b, rando.NewSource, 100) }
func BenchmarkRandoSrc_1000(b *testing.B)  { benchsrc(b, rando.NewSource, 1000) }
func BenchmarkRandoSrc_10000(b *testing.B) { benchsrc(b, rando.NewSource, 10000) }

func BenchmarkCryptoRandSrc_1(b *testing.B)     { benchsrc(b, crando.NewCrand, 1) }
func BenchmarkCryptoRandSrc_10(b *testing.B)    { benchsrc(b, crando.NewCrand, 10) }
func BenchmarkCryptoRandSrc_100(b *testing.B)   { benchsrc(b, crando.NewCrand, 100) }
func BenchmarkCryptoRandSrc_1000(b *testing.B)  { benchsrc(b, crando.NewCrand, 1000) }
func BenchmarkCryptoRandSrc_10000(b *testing.B) { benchsrc(b, crando.NewCrand, 10000) }

func BenchmarkMathRandGlobalSrc_1(b *testing.B)     { benchsrc(b, NewStdGlob, 1) }
func BenchmarkMathRandGlobalSrc_10(b *testing.B)    { benchsrc(b, NewStdGlob, 10) }
func BenchmarkMathRandGlobalSrc_100(b *testing.B)   { benchsrc(b, NewStdGlob, 100) }
func BenchmarkMathRandGlobalSrc_1000(b *testing.B)  { benchsrc(b, NewStdGlob, 1000) }
func BenchmarkMathRandGlobalSrc_10000(b *testing.B) { benchsrc(b, NewStdGlob, 10000) }

func BenchmarkRandoGlobalSrc_1(b *testing.B)     { benchsrc(b, rando.Pooled, 1) }
func BenchmarkRandoGlobalSrc_10(b *testing.B)    { benchsrc(b, rando.Pooled, 10) }
func BenchmarkRandoGlobalSrc_100(b *testing.B)   { benchsrc(b, rando.Pooled, 100) }
func BenchmarkRandoGlobalSrc_1000(b *testing.B)  { benchsrc(b, rando.Pooled, 1000) }
func BenchmarkRandoGlobalSrc_10000(b *testing.B) { benchsrc(b, rando.Pooled, 10000) }

func BenchmarkCryptoGlobalSrc_1(b *testing.B)     { benchsrc(b, crando.Pooled, 1) }
func BenchmarkCryptoGlobalSrc_10(b *testing.B)    { benchsrc(b, crando.Pooled, 10) }
func BenchmarkCryptoGlobalSrc_100(b *testing.B)   { benchsrc(b, crando.Pooled, 100) }
func BenchmarkCryptoGlobalSrc_1000(b *testing.B)  { benchsrc(b, crando.Pooled, 1000) }
func BenchmarkCryptoGlobalSrc_10000(b *testing.B) { benchsrc(b, crando.Pooled, 10000) }
