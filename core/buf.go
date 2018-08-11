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

// Sample buffer
type Buf [AUDIO_BUFSIZE]float32

func (a *Buf) Mul(b *Buf) {
	for i := 0; i < AUDIO_BUFSIZE; i++ {
		a[i] *= b[i]
	}
}

func (a *Buf) Add(b *Buf) {
	for i := 0; i < AUDIO_BUFSIZE; i++ {
		a[i] += b[i]
	}
}

func (a *Buf) MulScalar(k float32) {
	for i := 0; i < AUDIO_BUFSIZE; i++ {
		a[i] *= k
	}
}

func (a *Buf) AddScalar(k float32) {
	for i := 0; i < AUDIO_BUFSIZE; i++ {
		a[i] += k
	}
}

func (a *Buf) Copy(b *Buf) {
	for i := 0; i < AUDIO_BUFSIZE; i++ {
		a[i] = b[i]
	}
}

func (a *Buf) Zero() {
	for i := 0; i < AUDIO_BUFSIZE; i++ {
		a[i] = 0
	}
}

func (a *Buf) Set(k float32) {
	for i := 0; i < AUDIO_BUFSIZE; i++ {
		a[i] = k
	}
}

func (a *Buf) Equals(b *Buf) bool {
	for i := 0; i < AUDIO_BUFSIZE; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func (a *Buf) String() string {
	s := make([]string, AUDIO_BUFSIZE)
	for i := 0; i < AUDIO_BUFSIZE; i++ {
		s[i] = fmt.Sprintf("%.3f", a[i])
	}
	return fmt.Sprintf("[%s]", strings.Join(s, ","))
}

//-----------------------------------------------------------------------------
