//-----------------------------------------------------------------------------
/*

Random Functions

*/
//-----------------------------------------------------------------------------

package babi

import "math"

//-----------------------------------------------------------------------------

type Rand struct {
	state uint32
}

func NewRand(seed uint32) *Rand {
	if seed == 0 {
		seed = 1
	}
	return &Rand{seed}
}

// Return a random uint32_t (0..0x7fffffff)
func (r *Rand) Uint32() uint32 {
	r.state = ((r.state * 1103515245) + 12345) & 0x7fffffff
	return r.state
}

// Return a random float from -1..1
func (r *Rand) Float() float32 {
	i := r.Uint32()
	i |= (i << 1) & 0x80000000
	i = (i & 0x807fffff) | (126 << 23)
	return math.Float32frombits(i)
}

//-----------------------------------------------------------------------------
