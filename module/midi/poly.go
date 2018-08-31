//-----------------------------------------------------------------------------
/*

Polyphonic Module

Manage concurrent instances (voices) of a given sub-module.

TODO currently has a single output and assumes the sub-modules
have a single output.

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
		Name: "polyphonic",
		In: []core.PortInfo{
			{"midi_in", "midi input", core.PortType_EventMIDI, 0},
		},
		Out: []core.PortInfo{
			{"out", "output", core.PortType_AudioBuffer, 0},
		},
	}
}

//-----------------------------------------------------------------------------

type voiceInfo struct {
	note   uint8       // midi note value
	module core.Module // voice module
}

type polyModule struct {
	submodule      func() core.Module // new function for voice sub-module
	voice          []voiceInfo        // voices
	idx            int                // round-robin index for voice slice
	bend           float32            // pitch bending value (for all voices)
	frequency_ctrl uint               // sub-module frequency control id
	gate_ctrl      uint               // sub-module gate control id
}

func NewPoly(sm func() core.Module, maxvoices uint) core.Module {
	log.Info.Printf("")
	return &polyModule{
		submodule:      sm,
		voice:          make([]voiceInfo, maxvoices),
		frequency_ctrl: sm().Info().GetPortID("frequency"),
		gate_ctrl:      sm().Info().GetPortID("gate"),
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
	v.module.Event(core.NewEventFloat(m.frequency_ctrl, f))
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
				// note: vel=0 is the same as note off (gate=0).
				v.module.Event(core.NewEventFloat(m.gate_ctrl, vel))
			} else {
				if vel != 0 {
					v := m.voiceAlloc(me.GetNote())
					if v != nil {
						v.module.Event(core.NewEventFloat(m.gate_ctrl, vel))
					} else {
						log.Info.Printf("unable to allocate new voice")
					}
				}
			}
		case core.EventMIDI_NoteOff:
			v := m.voiceLookup(me.GetNote())
			if v != nil {
				// send a note off control event
				// ignoring the note off velocity (for now)
				v.module.Event(core.NewEventFloat(m.gate_ctrl, 0))
			}
		case core.EventMIDI_PitchWheel:
			// get the pitch bend value
			m.bend = core.MIDI_PitchBend(me.GetPitchWheel())
			// update all active voices
			for i := range m.voice {
				v := &m.voice[i]
				if v.module != nil {
					f := core.MIDI_ToFrequency(float32(v.note) + m.bend)
					v.module.Event(core.NewEventFloat(m.frequency_ctrl, f))
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
	out := buf[0]
	var vout core.Buf
	// run each voice
	for i := range m.voice {
		vm := m.voice[i].module
		if vm != nil && vm.Active() {
			// get the voice output
			vm.Process(&vout)
			// accumulate in the output buffer
			out.Add(&vout)
		}
	}
}

// Active return true if the module has non-zero output.
func (m *polyModule) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
