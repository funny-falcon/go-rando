# go-rando
Practially fast and practically secure random generator for Go.

Technically, it is entropy "multiplicator" that takes crypto/rand as source
of entropy and then applies high quality permutation.

To produce each 128bit of output generator:
- xor state with:
  - 64bit per generator counter
  - 64bit per generator unique id
  - 64bit salt randomly taken from global entropy pool of 128 elements.
- performs 5 rounds of SipHash permutaion
- takes two elements and xor it with per-generator mask

Generator is initialized in following way:
- id is get by permuting incremented global value.
  Incrementing global value and using invertable permutation, it quaranteed,
  that every generator in this process has unique id.
- counter is initialized with global random value xor-ed with id
- small generators, that produces indices to entropy pool, is derived from id in
  non-linear way.
- permutaion state and per-generator mask is then randomly taken from entropy pool.

Entropy pool of 128 uint64 is refreshed 10 times a second from crypto/rand.
All other global variables are also initialized from crypto/rand at process start.

# Installation
```sh
go get github.com/funny-falcon/go-rando
```

# Usage
```go
import "fmt"
import "math/rand"
import rando "github.com/funny-falcon/go-rando"

func main() {
	// global methods are same to math/rand
	fmt.Printf("%x %d\n", rando.Uint64(), rando.Intn(1000))
	// for generating many values it is better to wrap it to rand.New
	// Note: it is much cheaper to call rando.NewSource() than rand.NewSource()
	srng := rand.New(rando.NewSource())
	fmt.Printf("%x %d\n", srng.Uint64(), srng.Intn(1000))
}
```

# Quality

Generator uses permutation from cryptographic function, and it avoids presumable
pitfal of that permutation, so it is expected to pass any statictical tests.
It seems that the expectations are justified: dieharder, practrand and BigCrush
were completed without any fail. (results are in test.out folder)

# Performance

Despite being practically secure, Rando is almost as fast as math/rand:
- it is much smaller to allocate and initialize, because default math/rand
  source consists of 637 64bit numbers, and they all have to be initialized
  on seeding.
- it is just ~10% slower in long term. I don't know really why, may be still
  because of smaller footprint, and because it does less slice accesses.

You may see benchmark results on Intel i5-5200U CPU @ 2.20GHz in bench.out
and benchsrc.out
- bench.out - is benchmark of different math/rand.Source implementaions wrapped
  into math/rand.Rand, and math.Int32n method were called.
- benchsrc.out - is benchmark of different math/rand.Source64 implementations
  without wrapping into math/rand.Rand.

Also it was tested against libsodium chacha20 (benchmark folder). Go's rando implementation
is just 2.1x slower than ChaCha20. Implemented in C, rando is just 1.25x slower.
If number of iterations in permutation is decreased to 4, C rando is as fast as
ChaCha20.
(But, to be honestly, ChaCha8 is not yet broken, and it will be twice faster.)
