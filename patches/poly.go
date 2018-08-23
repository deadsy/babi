//-----------------------------------------------------------------------------
/*

Polyphonic Patch

Manage multiple instances of a given subpatch.

*/
//-----------------------------------------------------------------------------

package patches

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/log"
)

//-----------------------------------------------------------------------------

const MAX_VOICES = 16

//-----------------------------------------------------------------------------

type voiceInfo struct {
	note  uint8 // midi note value
	patch core.Patch
}

type polyPatch struct {
	newvoice func() core.Patch     // new function for voice subpatch
	idx      uint                  // round-robin index for voice array
	voice    [MAX_VOICES]voiceInfo // voice array
}

func NewPolyPatch(newvoice func() core.Patch) core.Patch {
	log.Info.Printf("")
	return &polyPatch{
		newvoice: newvoice,
	}
}

//-----------------------------------------------------------------------------

// Run the patch.
func (p *polyPatch) Process(in, out []*core.Buf) {
	for i := range p.voice {
		vp := p.voice[i].patch
		if vp != nil && vp.Active() {
			vp.Process(nil, []*core.Buf{out[0]})
		}
	}
}

// Process a patch event.
func (p *polyPatch) Event(e *core.Event) {
	switch e.GetType() {
	case core.Event_MIDI:
		me := e.GetMIDIEvent()
		switch me.GetType() {
		case core.MIDIEvent_NoteOn:
			v := p.voiceLookup(me.GetNote())
			vel := core.MIDI_Map(me.GetVelocity(), 0, 1)
			if v != nil {
				if vel == 0 {
					// velocity 0 == note off
					v.patch.Event(core.NewCtrlEvent(core.CtrlEvent_NoteOff, vel))
				} else {
					// trigger the note again
					v.patch.Event(core.NewCtrlEvent(core.CtrlEvent_NoteOn, vel))
				}
			} else {
				if vel != 0 {
					v := p.voiceAlloc(me.GetNote())
					if v != nil {
						v.patch.Event(core.NewCtrlEvent(core.CtrlEvent_NoteOn, vel))
					} else {
						log.Info.Printf("unable to allocate new voice")
					}
				}
			}
		case core.MIDIEvent_NoteOff:
			v := p.voiceLookup(me.GetNote())
			if v != nil {
				// send a note off control event
				vel := core.MIDI_Map(me.GetVelocity(), 0, 1)
				v.patch.Event(core.NewCtrlEvent(core.CtrlEvent_NoteOff, vel))
			}
		default:
			log.Info.Printf("unhandled midi event type %s", me)

		}
	default:
		log.Info.Printf("unhandled event type %s", e)
	}
}

// Return true if the patch has non-zero output.
func (p *polyPatch) Active() bool {
	return true
}

func (p *polyPatch) Stop() {
	// do nothing
}

//-----------------------------------------------------------------------------

// lookup a voice by note
func (p *polyPatch) voiceLookup(note uint8) *voiceInfo {
	for i := range p.voice {
		if p.voice[i].patch != nil && p.voice[i].note == note {
			return &p.voice[i]
		}
	}
	return nil
}

// allocate a new voice
func (p *polyPatch) voiceAlloc(note uint8) *voiceInfo {
	// Currently doing simple round robin allocation.
	v := &p.voice[p.idx]
	p.idx += 1
	if p.idx == MAX_VOICES {
		p.idx = 0
	}
	// stop an existing patch on this voice
	if v.patch != nil {
		v.patch.Stop()
	}
	// setup the new voice
	v.note = note
	v.patch = p.newvoice()
	return v
}

//-----------------------------------------------------------------------------
