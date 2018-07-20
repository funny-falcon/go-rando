package crando

import (
	crand "crypto/rand"
	"encoding/binary"
	"io"
	"math/rand"
)

// CrandReader - implementation of io.Reader that uses crypto/rand.Read
type CrandReader struct{}

// Read - implementation of io.Reader
func (c CrandReader) Read(buf []byte) (int, error) {
	return crand.Read(buf)
}

// Crand - math/rand.Source implementation, that uses crypto/rand.Read
// Note: it is not goroutine safe. Serialize access in your way.
type Crand struct {
	pos int
	buf [2040]byte
}

// NewCrand returns new Crand
// You can always use &Crand{} or new(Crand). It is just for convenience.
func NewCrand() rand.Source {
	return &Crand{}
}

// ensure there is enough data
func (c *Crand) ensure() {
	if c.pos < 8 {
		_, err := io.ReadFull(CrandReader{}, c.buf[:])
		if err != nil {
			panic(err)
		}
		c.pos = len(c.buf)
	}
}

// Uint64 returns random uint64 - implementation of rand.Source64
func (c *Crand) Uint64() uint64 {
	c.ensure()
	c.pos -= 8
	return binary.LittleEndian.Uint64(c.buf[c.pos:])
}

// Int63 returns random int63 - implementation of rand.Source
func (c *Crand) Int63() int64 {
	return int64(c.Uint64() >> 1)
}

// Seed - fake implementation of rand.Source.Reed
// It is meaningless to seed secure random generator.
func (c *Crand) Seed(i int64) {}
