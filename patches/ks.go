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
	ks *osc.KarplusStrong
}

func NewKSPatch() core.Patch {
	log.Info.Printf("")
	return &ksPatch{
		ks: osc.NewKarplusStrong(),
	}
}

//-----------------------------------------------------------------------------

// Run the patch.
func (p *ksPatch) Process(in, out []*core.Buf) {
	// generate the ks wave
	p.ks.Process(out[0])
}

// Process a patch event.
func (p *ksPatch) Event(e *core.Event) {
	log.Info.Printf("event %s", e)
	switch e.GetType() {
	case core.Event_Ctrl:
		ce := e.GetCtrlEvent()
		switch ce.GetType() {
		case core.CtrlEvent_NoteOn:
			p.ks.Pluck() // velocity?
		case core.CtrlEvent_NoteOff:
			// ignore
		case core.CtrlEvent_Frequency:
			p.ks.SetFrequency(ce.GetVal())
		case core.CtrlEvent_Attenuate:
			p.ks.SetAttenuate(ce.GetVal())
		default:
			log.Info.Printf("unhandled ctrl event %s", ce)
		}
	default:
		log.Info.Printf("unhandled event %s", e)
	}
}

// Return true if the patch has non-zero output.
func (p *ksPatch) Active() bool {
	return true
}

func (p *ksPatch) Stop() {
	log.Info.Printf("")
	// do nothing
}

//-----------------------------------------------------------------------------
