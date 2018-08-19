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
func (p *simplePatch) Process() {
	var env, x core.Buf
	// generate the envelope
	p.adsr.Process(&env)
	// generate the sine wave
	p.sine.Process(&x)
	// apply the envelope
	x.Mul(&env)
	// output
	p.parent.Out(&x)
}

// Process a patch event.
func (p *simplePatch) Event(e *core.Event) {
}

// Return true if the patch has non-zero output.
func (p *simplePatch) Active() bool {
	return p.adsr.Active()
}

// Output to the parent patch.
func (p *simplePatch) Out(out ...*core.Buf) {
}

//-----------------------------------------------------------------------------

/*

var Simple = core.PatchInfo{
	Name: "simple",
	New:  NewPatch,
}

//-----------------------------------------------------------------------------

// voice state
type voice struct {
	note uint      // note number
	p    *Patch    // parent patch
	adsr *env.ADSR // adsr envelope
	sine *osc.Sine // sine oscillator
}

// Return true if the voice has non-zero output.
func (v *voice) active() bool {
	return v.adsr.Active()
}

func (v *voice) setFrequency() {
	v.ks.SetFrequency(core.MIDI_ToFrequency(float32(v.note) + v.p.bend))
}

//-----------------------------------------------------------------------------

// patch state
type Patch struct {
	s     *core.Synth  // parent synth
	voice []*voice     // set of voices on this patch
	pan   *audio.Pan   // left right panning
	out   *audio.OutLR //audio output
	bend  float32      // pitch bend
}

// NewPatch returns a new patch (with no voices) for the synth.
func NewPatch(s *core.Synth) core.Patch {
	p := &Patch{
		s:   s,
		pan: audio.NewPan(),
		out: audio.NewOutLR(s),
	}
	p.pan.SetPan(0.5)
	p.pan.SetVol(1.0)
	return p
}

// StartVoice allocates and starts a new voice for the patch.
func (p *Patch) StartVoice(note uint) {
	v := &voice{
		note: note,
		p:    p,
		adsr: env.NewAD(0.1, 1.0),
		sine: osc.NewSine(),
	}
	p.voice = append(p.voice, v)
	v.setFrequency()
}

// StopVoice de-allocates and stops voices on the patch.
func (p *Patch) StopVoice(note uint) {
	for i, v := range p.voice {
		if v != nil && v.note == note {
			p.voice[i] = nil
		}
	}
}

func (p *Patch) Process() {
	var out core.Buf
	for _, v := range p.voice {
		if v != nil && v.active() {
			var env, x core.Buf
			// generate the envelope
			v.adsr.Process(&env)
			// generate the sine wave
			v.sine.Process(&x)
			// apply the envelope
			x.Mul(&env)
			// add to the output
			out.Add(&x)
		}
	}
	// pan to left/right channels
	var out_l, out_r core.Buf
	p.pan.Process(&out, &out_l, &out_r)
	// stereo output
	p.out.Process(&out_l, &out_r)
}

*/

//-----------------------------------------------------------------------------
