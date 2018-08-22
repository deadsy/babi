//-----------------------------------------------------------------------------
/*

A simple patch - an AD envelope on a sine wave.

*/
//-----------------------------------------------------------------------------

package patches

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/log"
	"github.com/deadsy/babi/objects/env"
	"github.com/deadsy/babi/objects/osc"
)

//-----------------------------------------------------------------------------

type simplePatch struct {
	parent core.Patch
	adsr   *env.ADSR // adsr envelope
	sine   *osc.Sine // sine oscillator
}

func NewSimplePatch() core.Patch {
	log.Info.Printf("")
	return &simplePatch{}
}

//-----------------------------------------------------------------------------

// Run the patch.
func (p *simplePatch) Process(in, out []*core.Buf) {
	var env core.Buf
	// generate the envelope
	p.adsr.Process(&env)
	// generate the sine wave
	p.sine.Process(out[0])
	// apply the envelope
	out[0].Mul(&env)
}

// Process a patch event.
func (p *simplePatch) Event(e *core.Event) {
}

// Return true if the patch has non-zero output.
func (p *simplePatch) Active() bool {
	return p.adsr.Active()
}

func (p *simplePatch) Stop() {
	// do nothing
}

//-----------------------------------------------------------------------------
