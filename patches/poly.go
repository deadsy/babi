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
	note  uint
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
	switch e.Etype {
	case core.Event_MIDI:
		p.midiEvent(e.Info.(*core.MIDIEvent))
	default:
		log.Info.Printf("unhandled event type %s", e)
	}
}

// Return true if the patch has non-zero output.
func (p *polyPatch) Active() bool {
	return true
}

//-----------------------------------------------------------------------------

// Handle a MIDI event.
func (p *polyPatch) midiEvent(e *core.MIDIEvent) {
	switch e.Etype {
	case core.MIDIEvent_NoteOn:
		note := e.GetNote()
		vel := e.GetVelocity()
		if vel == 0 {
			// velocity 0 == note off
			p.noteOff(note, vel)
		} else {
			p.noteOn(note, vel)
		}
	case core.MIDIEvent_NoteOff:
		note := e.GetNote()
		vel := e.GetVelocity()
		p.noteOff(note, vel)
	default:
		log.Info.Printf("unhandled midi event type %s", e)
	}
}

//-----------------------------------------------------------------------------

// lookup a voice by note
func (p *polyPatch) voiceLookup(note uint) *voiceInfo {
	for i := range p.voice {
		if p.voice[i].patch != nil && p.voice[i].note == note {
			return &p.voice[i]
		}
	}
	return nil
}

//-----------------------------------------------------------------------------

// Handle a note on event.
func (p *polyPatch) noteOn(note, vel uint) {
	log.Info.Printf("note %d vel %d", note, vel)
}

// Handle a note off event.
func (p *polyPatch) noteOff(note, vel uint) {
	log.Info.Printf("note %d vel %d", note, vel)
}

//-----------------------------------------------------------------------------
