//-----------------------------------------------------------------------------
/*

Sine Oscillator

*/
//-----------------------------------------------------------------------------

package osc

import "github.com/deadsy/babi/core"

//-----------------------------------------------------------------------------

type Sine struct {
}

func NewSine() *Sine {
	return &Sine{}
}

func (o *Sine) Process(out *core.Buf) {
}

//-----------------------------------------------------------------------------
