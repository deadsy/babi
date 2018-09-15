//-----------------------------------------------------------------------------
/*

Buffer types/operations

*/
//-----------------------------------------------------------------------------

package core

import (
	"fmt"
	"strings"
)

//-----------------------------------------------------------------------------
// Sample Buffers (at audio sample rate)

// Buf is an audio sample buffer.
type Buf [AudioBufferSize]float32

// Mul multiplies two buffers, a := a * b
func (a *Buf) Mul(b *Buf) {
	for i := 0; i < AudioBufferSize; i++ {
		a[i] *= b[i]
	}
}

// Add adds two buffers, a := a + b
func (a *Buf) Add(b *Buf) {
	for i := 0; i < AudioBufferSize; i++ {
		a[i] += b[i]
	}
}

// MulScalar multiplies a buffer by a scalar, a := [k * a0, k * a1, ...]
func (a *Buf) MulScalar(k float32) {
	for i := 0; i < AudioBufferSize; i++ {
		a[i] *= k
	}
}

// AddScalar adds a scalar to a buffer, a := [k + a0, k + a1, ...]
func (a *Buf) AddScalar(k float32) {
	for i := 0; i < AudioBufferSize; i++ {
		a[i] += k
	}
}

// Copy copies a buffer, a := b
func (a *Buf) Copy(b *Buf) {
	for i := 0; i < AudioBufferSize; i++ {
		a[i] = b[i]
	}
}

// Zero zeroes a buffer, a := [0, 0, ...]
func (a *Buf) Zero() {
	for i := 0; i < AudioBufferSize; i++ {
		a[i] = 0
	}
}

// Set sets a buffer to a value, a := [k, k, ...]
func (a *Buf) Set(k float32) {
	for i := 0; i < AudioBufferSize; i++ {
		a[i] = k
	}
}

// Equal tests if two buffers are equal, returns true if a == b.
func (a *Buf) Equal(b *Buf) bool {
	for i := 0; i < AudioBufferSize; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// String returns a string representation of the buffer.
func (a *Buf) String() string {
	s := make([]string, AudioBufferSize)
	for i := 0; i < AudioBufferSize; i++ {
		s[i] = fmt.Sprintf("%.3f", a[i])
	}
	return fmt.Sprintf("[%s]", strings.Join(s, ","))
}

//-----------------------------------------------------------------------------
