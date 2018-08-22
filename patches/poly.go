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
		p.midiEvent(e)
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
func (p *polyPatch) voiceLookup(note uint) *voiceInfo {
	for i := range p.voice {
		if p.voice[i].patch != nil && p.voice[i].note == note {
			return &p.voice[i]
		}
	}
	return nil
}

// allocate a new voice
func (p *polyPatch) voiceAlloc(note uint) *voiceInfo {
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
// MIDI events

// Handle a MIDI note off event.
func (p *polyPatch) midiNoteOff(e *core.Event) {
	me := e.Info.(*core.MIDIEvent)
	note := me.GetNote()
	vel := me.GetVelocity()
	log.Info.Printf("note %d vel %d", note, vel)
	v := p.voiceLookup(note)
	if v != nil {
		v.patch.Event(e)
	}
}

// Handle a MIDI note on event.
func (p *polyPatch) midiNoteOn(e *core.Event) {
	me := e.Info.(*core.MIDIEvent)
	note := me.GetNote()
	vel := me.GetVelocity()
	log.Info.Printf("note %d vel %d", note, vel)
	if vel == 0 {
		// velocity 0 == note off
    me.EType = core.MIDIEvent_NoteOff
    p.midiNoteOff(e)
	}
	v := p.voiceLookup(note)
	if v == nil {
		v = p.voiceAlloc(note)
	}
	if v != nil {
		v.patch.Event(e)
	}
}

// Handle a MIDI event.
func (p *polyPatch) midiEvent(e *core.Event) {
	me := e.Info.(*core.MIDIEvent)
	switch me.Etype {
	case core.MIDIEvent_NoteOn:
		p.midiNoteOff(e)
	case core.MIDIEvent_NoteOff:
		p.midiNoteOn(e)
	default:
		log.Info.Printf("unhandled midi event type %s", e)
	}
}

//-----------------------------------------------------------------------------
