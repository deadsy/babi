//-----------------------------------------------------------------------------
/*

Buffer types/operations

*/
//-----------------------------------------------------------------------------

package core

//-----------------------------------------------------------------------------
// Sample Buffers (at audio sample rate)

const SAMPLES_PER_SBUF = 64

// S (sample) buffer
type SBuf [SAMPLES_PER_SBUF]float32

// Set_SK sets a sample buffer to a constant value.
func Set_SK(a *SBuf, k float32) {
	for i := 0; i < SAMPLES_PER_SBUF; i++ {
		a[i] = k
	}
}

// Mul_SS multiples two sample buffers.
func Mul_SS(a, b *SBuf) {
	for i := 0; i < SAMPLES_PER_SBUF; i++ {
		a[i] *= b[i]
	}
}

// Mul_SK multiples a sample buffer by a scalar.
func Mul_SK(a *SBuf, k float32) {
	for i := 0; i < SAMPLES_PER_SBUF; i++ {
		a[i] *= k
	}
}

// Add_SS adds two sample buffers.
func Add_SS(a, b *SBuf) {
	for i := 0; i < SAMPLES_PER_SBUF; i++ {
		a[i] += b[i]
	}
}

// Add_SK multiples a sample buffer by a scalar.
func Add_SK(a *SBuf, k float32) {
	for i := 0; i < SAMPLES_PER_SBUF; i++ {
		a[i] += k
	}
}

// Copy_S copies a sample buffer.
func Copy_S(dst, src *SBuf) {
	dst = src
}

// Copy_SK copies a sample buffer and multiplies by K.
func Copy_SK(dst, src *SBuf, k float32) {
	Copy_S(dst, src)
	Mul_SK(dst, k)
}

// Equal_SS tests the equality of two sample buffers.
func Equal_SS(a, b *SBuf) bool {
	for i := 0; i < SAMPLES_PER_SBUF; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

//-----------------------------------------------------------------------------
// Event buffers (slow moving)

const SAMPLES_PER_EBUF = 8

// E (event) buffer
type EBuf [SAMPLES_PER_EBUF]float32

//-----------------------------------------------------------------------------
