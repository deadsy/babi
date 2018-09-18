//-----------------------------------------------------------------------------
/*

Random Functions

*/
//-----------------------------------------------------------------------------

package core

import "math"

//-----------------------------------------------------------------------------

// Rand contains the state for a simple PRNG.
type Rand struct {
	state uint32
}

// NewRand returns a simple PRNG object.
func NewRand(seed uint32) *Rand {
	if seed == 0 {
		seed = 1
	}
	return &Rand{seed}
}

// Uint32 returns a random uint32_t (0..0x7fffffff)
func (r *Rand) Uint32() uint32 {
	r.state = ((r.state * 1103515245) + 12345) & 0x7fffffff
	return r.state
}

// Float returns a random float from -1..1
func (r *Rand) Float() float32 {
	i := (r.Uint32() & 0x007fffff) | (128 << 23)
	return math.Float32frombits(i) - 3
}

//-----------------------------------------------------------------------------
