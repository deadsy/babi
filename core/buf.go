//-----------------------------------------------------------------------------
/*

Buffer types/operations

*/
//-----------------------------------------------------------------------------

package core

//-----------------------------------------------------------------------------
// Sample Buffers (at audio sample rate)

const SAMPLES_PER_BUF = 64

// Sample buffer
type Buf [SAMPLES_PER_BUF]float32

// Return a new buffer initialised to k
func NewBuf(k float32) *Buf {
	var b Buf
	for i := 0; i < SAMPLES_PER_BUF; i++ {
		b[i] = k
	}
	return &b
}

func (a *Buf) Mul(b *Buf) {
	for i := 0; i < SAMPLES_PER_BUF; i++ {
		a[i] *= b[i]
	}
}

func (a *Buf) Add(b *Buf) {
	for i := 0; i < SAMPLES_PER_BUF; i++ {
		a[i] += b[i]
	}
}

func (a *Buf) MulScalar(k float32) {
	for i := 0; i < SAMPLES_PER_BUF; i++ {
		a[i] *= k
	}
}

func (a *Buf) AddScalar(k float32) {
	for i := 0; i < SAMPLES_PER_BUF; i++ {
		a[i] += k
	}
}

func (a *Buf) Copy(b *Buf) {
	a = b
}

func (a *Buf) Equals(b *Buf) bool {
	for i := 0; i < SAMPLES_PER_BUF; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

//-----------------------------------------------------------------------------
