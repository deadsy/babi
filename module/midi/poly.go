//-----------------------------------------------------------------------------
/*

Polyphonic Module

Manage concurrent instances (voices) of a given sub-module.

*/
//-----------------------------------------------------------------------------

package midi

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/log"
)

//-----------------------------------------------------------------------------

// Info returns the module information.
func (m *polyModule) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		In: []core.PortInfo{
			{"midi", "midi input", core.PortType_EventMIDI},
		},
		Out: nil,
	}
}

//-----------------------------------------------------------------------------

type voiceInfo struct {
	note   uint8       // midi note value
	module core.Module // voice module
}

type polyModule struct {
	submodule    func() core.Module // new function for voice sub-module
	voice        []voiceInfo        // voices
	idx          int                // round-robin index for voice slice
	bend         float32            // pitch bending value (for all voices)
	frequency_id uint               // sub-module control id
	noteon_id    uint               // sub-module control id
	noteoff_id   uint               // sub-module control id
}

func NewPoly(sm func() core.Module, maxvoices uint) core.Module {
	log.Info.Printf("")
	return &polyModule{
		submodule: sm,
		voice:     make([]voiceInfo, maxvoices),
	}
}

// Stop performs any cleanup of a module.
func (m *polyModule) Stop() {
	log.Info.Printf("")
}

//-----------------------------------------------------------------------------
// Events

// voiceLookup returns the voice for this MIDI note (or nil).
func (m *polyModule) voiceLookup(note uint8) *voiceInfo {
	for i := range m.voice {
		if m.voice[i].module != nil && m.voice[i].note == note {
			return &m.voice[i]
		}
	}
	return nil
}

// voiceAlloc allocates a new subpatch voice for a MIDI note.
func (m *polyModule) voiceAlloc(note uint8) *voiceInfo {
	log.Info.Printf("")
	// Currently doing simple round robin allocation.
	v := &m.voice[m.idx]
	m.idx += 1
	if m.idx == len(m.voice) {
		m.idx = 0
	}
	// stop an existing patch on this voice
	if v.module != nil {
		v.module.Stop()
	}
	// setup the new voice
	v.note = note
	v.module = m.submodule()
	// set the voice frequency
	f := core.MIDI_ToFrequency(float32(v.note) + m.bend)
	v.module.Event(core.NewEventFloat(m.frequency_id, f))
	return v
}

// Event processes a module event.
func (m *polyModule) Event(e *core.Event) {
	log.Info.Printf("event %s", e)
	switch e.GetType() {
	case core.Event_MIDI:
		me := e.GetEventMIDI()
		switch me.GetType() {
		case core.EventMIDI_NoteOn:
			v := m.voiceLookup(me.GetNote())
			vel := core.MIDI_Map(me.GetVelocity(), 0, 1)
			if v != nil {
				if vel == 0 {
					// velocity 0 == note off
					v.module.Event(core.NewEventFloat(m.noteoff_id, vel))
				} else {
					// trigger the note again
					v.module.Event(core.NewEventFloat(m.noteon_id, vel))
				}
			} else {
				if vel != 0 {
					v := m.voiceAlloc(me.GetNote())
					if v != nil {
						v.module.Event(core.NewEventFloat(m.noteon_id, vel))
					} else {
						log.Info.Printf("unable to allocate new voice")
					}
				}
			}
		case core.EventMIDI_NoteOff:
			v := m.voiceLookup(me.GetNote())
			if v != nil {
				// send a note off control event
				vel := core.MIDI_Map(me.GetVelocity(), 0, 1)
				v.module.Event(core.NewEventFloat(m.noteoff_id, vel))
			}
		case core.EventMIDI_PitchWheel:
			// get the pitch bend value
			m.bend = core.MIDI_PitchBend(me.GetPitchWheel())
			// update all active voices
			for i := range m.voice {
				v := &m.voice[i]
				if v.module != nil {
					f := core.MIDI_ToFrequency(float32(v.note) + m.bend)
					v.module.Event(core.NewEventFloat(m.frequency_id, f))
				}
			}
		default:
			log.Info.Printf("unhandled midi event %s", me)
		}
	default:
		log.Info.Printf("unhandled event %s", e)
	}
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *polyModule) Process(buf ...*core.Buf) {
	for i := range m.voice {
		vm := m.voice[i].module
		if vm != nil && vm.Active() {
			vm.Process(buf...)
		}
	}
}

// Active return true if the module has non-zero output.
func (m *polyModule) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
