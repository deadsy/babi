//-----------------------------------------------------------------------------
/*

Karplus Strong Patch

*/
//-----------------------------------------------------------------------------

package patches

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/log"
	"github.com/deadsy/babi/objects/audio"
	"github.com/deadsy/babi/objects/osc"
)

//-----------------------------------------------------------------------------

var KarplusStrongInfo = core.PatchInfo{
	Name: "karplus_strong",
	New:  NewKarplusStrong,
}

type KarplusStrong struct {
	ks  *osc.KarplusStrong
	pan *audio.Pan
	out *audio.OutLR
}

func NewKarplusStrong(s *core.Synth) core.Patch {
	log.Info.Printf("")

	p := &KarplusStrong{
		ks:  osc.NewKarplusStrong(),
		pan: audio.NewPan(),
		out: audio.NewOutLR(s),
	}

	p.ks.SetFrequency(440.0)
	p.ks.SetAttenuate(1.0)
	p.ks.Pluck()
	p.pan.SetPan(0.5)
	p.pan.SetVol(1.0)

	return p
}

func (p *KarplusStrong) Stop() {
	log.Info.Printf("")
}

func (p *KarplusStrong) Active() bool {
	return true
}

func (p *KarplusStrong) Process() {
	var out, out_l, out_r core.Buf
	// generate the ks wave
	p.ks.Process(&out)
	// pan to left/right channels
	p.pan.Process(&out, &out_l, &out_r)
	// stereo output
	p.out.Process(&out_l, &out_r)
}

//-----------------------------------------------------------------------------
