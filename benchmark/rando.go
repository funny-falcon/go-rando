package main

import (
	"encoding/binary"
	"math/rand"
	"os"

	rando "github.com/funny-falcon/go-rando"
)

func main() {
	rng := rando.NewSource().(rand.Source64)
	b := make([]byte, 1024*1024)
	for {
		for i := 0; i < 1024*1024; i += 8 {
			u := rng.Uint64()
			binary.LittleEndian.PutUint64(b[i:], u)
		}
		os.Stdout.Write(b)
	}
}
