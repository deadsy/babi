//-----------------------------------------------------------------------------
/*

Random Functions

See:
https://en.wikipedia.org/wiki/Linear_congruential_generator

*/
//-----------------------------------------------------------------------------

package core

import "math"

//-----------------------------------------------------------------------------

// Rand32 is a simple LCG PRNG.
type Rand32 struct {
	state uint32
}

// NewRand32 returns the state for a 32-bit PRNG.
func NewRand32(seed uint32) *Rand32 {
	if seed == 0 {
		seed = 1
	}
	return &Rand32{seed}
}

// Uint32 returns a random uint32.
func (r *Rand32) Uint32() uint32 {
	r.state = (r.state * 214013) + 2531011
	return r.state
}

// Float32 returns a random float32 from -1..1
func (r *Rand32) Float32() float32 {
	i := (r.Uint32() & 0x007fffff) | (128 << 23) // 2..4
	return math.Float32frombits(i) - 3           // -1..1
}

//-----------------------------------------------------------------------------

// Rand64 is a simple LCG PRNG.
type Rand64 struct {
	state uint64
}

// NewRand64 returns the state for a 64-bit PRNG.
func NewRand64(seed uint64) *Rand64 {
	if seed == 0 {
		seed = 1
	}
	return &Rand64{seed}
}

// Uint64 returns a random uint64.
func (r *Rand64) Uint64() uint64 {
	r.state = (r.state * 6364136223846793005) + 1442695040888963407
	return r.state
}

// Float64 returns a random float64 from -1..1
func (r *Rand64) Float64() float64 {
	i := (r.Uint64() & 0x000fffffffffffff) | (1024 << 52) // 2..4
	return math.Float64frombits(i) - 3                    // -1..1
}

//-----------------------------------------------------------------------------
