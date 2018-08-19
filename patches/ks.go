//-----------------------------------------------------------------------------
/*

Karplus Strong Patch

*/
//-----------------------------------------------------------------------------

package patches

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/log"
	"github.com/deadsy/babi/objects/osc"
)

//-----------------------------------------------------------------------------

type ksPatch struct {
	parent core.Patch
	ks     *osc.KarplusStrong
}

func NewKSPatch() core.Patch {
	log.Info.Printf("")
	return &ksPatch{}
}

//-----------------------------------------------------------------------------

// Run the patch.
func (p *ksPatch) Process() {
	var x core.Buf
	// generate the ks wave
	p.ks.Process(&x)
	// output
	p.parent.Out(&x)
}

// Process a patch event.
func (p *ksPatch) Event(e *core.Event) {
}

// Return true if the patch has non-zero output.
func (p *ksPatch) Active() bool {
	return true
}

// Output to the parent patch.
func (p *ksPatch) Out(out ...*core.Buf) {
}

//-----------------------------------------------------------------------------

/*


// voice state
type voice struct {
	note uint   // note number
	p    *Patch // parent patch
	ks   *osc.KarplusStrong
}

// Return true if the voice has non-zero output.
func (v *voice) active() bool {
	return true
}

func (v *voice) setAttenuate() {
	v.ks.SetAttenuate(v.p.attenuate)
}

func (v *voice) setFrequency() {
	v.ks.SetFrequency(core.MIDI_ToFrequency(float32(v.note) + v.p.bend))
}

//-----------------------------------------------------------------------------

type Patch struct {
	s         *core.Synth  // parent synth
	voice     []*voice     // set of voices on this patch
	pan       *audio.Pan   // left right panning
	out       *audio.OutLR //audio output
	bend      float32      // pitch bend
	attenuate float32      // karplus-strong attenuation
}

// NewPatch returns a new patch (with no voices) for the synth.
func NewPatch(s *core.Synth) core.Patch {
	p := &Patch{
		s:   s,
		pan: audio.NewPan(),
		out: audio.NewOutLR(s),
	}
	p.attenuate = 1.0
	p.pan.SetPan(0.5)
	p.pan.SetVol(1.0)
	return p
}

//-----------------------------------------------------------------------------

// StartVoice allocates and starts a new voice for the patch.
func (p *Patch) StartVoice(note uint) {
	v := &voice{
		note: note,
		p:    p,
		ks:   osc.NewKarplusStrong(),
	}
	p.voice = append(p.voice, v)
	v.setFrequency()
	v.setAttenuate()
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
			var x core.Buf
			// generate the ks wave
			v.ks.Process(&x)
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
