//-----------------------------------------------------------------------------
/*

Polyphonic Patch

Manage concurrent instances (voices) of a given subpatch.

*/
//-----------------------------------------------------------------------------

package patches

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/log"
)

//-----------------------------------------------------------------------------

type voiceInfo struct {
	note  uint8      // midi note value
	patch core.Patch // voice patch
}

type polyPatch struct {
	subpatch func() core.Patch // new function for voice subpatch
	voice    []voiceInfo       // voices
	idx      int               // round-robin index for voice slice
	bend     float32           // pitch bending value (for all voices)
}

func NewPolyPatch(subpatch func() core.Patch, maxvoices uint) core.Patch {
	log.Info.Printf("")
	return &polyPatch{
		subpatch: subpatch,
		voice:    make([]voiceInfo, maxvoices),
	}
}

//-----------------------------------------------------------------------------

// voiceLookup returns the voice for this MIDI note (or nil).
func (p *polyPatch) voiceLookup(note uint8) *voiceInfo {
	for i := range p.voice {
		if p.voice[i].patch != nil && p.voice[i].note == note {
			return &p.voice[i]
		}
	}
	return nil
}

// voiceAlloc allocates a new subpatch voice for a MIDI note.
func (p *polyPatch) voiceAlloc(note uint8) *voiceInfo {
	// Currently doing simple round robin allocation.
	v := &p.voice[p.idx]
	p.idx += 1
	if p.idx == len(p.voice) {
		p.idx = 0
	}
	// stop an existing patch on this voice
	if v.patch != nil {
		v.patch.Stop()
	}
	// setup the new voice
	v.note = note
	v.patch = p.subpatch()
	// set the voice frequency
	f := core.MIDI_ToFrequency(float32(v.note) + p.bend)
	v.patch.Event(core.NewCtrlEvent(core.CtrlEvent_Frequency, f))
	return v
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
		case core.MIDIEvent_PitchWheel:
			// get the pitch bend value
			p.bend = core.MIDI_PitchBend(me.GetPitchWheel())
			// update all active voices
			for i := range p.voice {
				v := &p.voice[i]
				if v.patch != nil {
					f := core.MIDI_ToFrequency(float32(v.note) + p.bend)
					v.patch.Event(core.NewCtrlEvent(core.CtrlEvent_Frequency, f))
				}
			}
		default:
			log.Info.Printf("unhandled midi event %s", me)
		}
	default:
		log.Info.Printf("unhandled event %s", e)
	}
}

// Return true if the patch has non-zero output.
func (p *polyPatch) Active() bool {
	return true
}

func (p *polyPatch) Stop() {
	log.Info.Printf("")
	// do nothing
}

//-----------------------------------------------------------------------------
