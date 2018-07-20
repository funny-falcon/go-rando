package rando

import (
	"sync/atomic"
	"time"

	"github.com/funny-falcon/go-rando/crando"
)

var crandg = crando.Crand{}

// id is a counter that guarantees that every generator instance has
// distinct per round mixin.
var genid = crandg.Uint64()

// some random increment and multiplicator used to permute
// per generator id.
// Sure increment is odd to have 2^64 unique identificators in one process.
// Sure multiplicator is odd to not loose bits.
var genidadd = crandg.Uint64() | 1
var genidmult = crandg.Uint64() | 0x200080001

// counter start offset.
// While it is similar for all generators, and it is not
// necessary for it to be random, lets make it random.
var genoffset = crandg.Uint64()

// counter for reseeder initialization
var genpeekcnt = crandg.Uint64()

// Seed space refreshes every 100ms from crypto/rand.Read and used as
// generator initializator and random mixin on every 128bit output.
var seedspace = [128]uint64{}

// refresh seedspace
func fillseed() {
	for i := range seedspace {
		atomic.StoreUint64(&seedspace[i], crandg.Uint64())
	}
}

func init() {
	fillseed()
	go func() {
		for range time.Tick(100 * time.Millisecond) {
			fillseed()
		}
	}()
}
