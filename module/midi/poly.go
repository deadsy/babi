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
		Name: "poly",
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
	ch        uint8              // MIDI channel
	submodule func() core.Module // new function for voice sub-module
	voice     []voiceInfo        // voices
	idx       int                // round-robin index for voice slice
	bend      float32            // pitch bending value (for all voices)
	freq      core.PortId        // sub-module frequency port id
	gate      core.PortId        // sub-module gate port id
}

// NewPoly returns a MIDI polyphonic voice control module.
func NewPoly(ch uint8, sm func() core.Module, maxvoices uint) core.Module {
	log.Info.Printf("")
	return &polyModule{
		ch:        ch,
		submodule: sm,
		voice:     make([]voiceInfo, maxvoices),
		freq:      sm().Info().GetPortId("frequency"),
		gate:      sm().Info().GetPortId("gate"),
	}
}

// Return the child modules.
func (m *polyModule) Child() []core.Module {
	var children []core.Module
	for i := range m.voice {
		if m.voice[i].module != nil {
			children = append(children, m.voice[i].module)
		}
	}
	return children
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
	core.SendEventFloatID(v.module, m.freq, f)
	return v
}

// Event processes a module event.
func (m *polyModule) Event(e *core.Event) {
	me := e.GetEventMIDIChannel(m.ch)
	if me != nil {
		switch me.GetType() {
		case core.EventMIDI_NoteOn:
			v := m.voiceLookup(me.GetNote())
			vel := core.MIDI_Map(me.GetVelocity(), 0, 1)
			if v != nil {
				// note: vel=0 is the same as note off (gate=0).
				core.SendEventFloatID(v.module, m.gate, vel)
			} else {
				if vel != 0 {
					v := m.voiceAlloc(me.GetNote())
					if v != nil {
						core.SendEventFloatID(v.module, m.gate, vel)
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
				core.SendEventFloatID(v.module, m.gate, 0)
			}
		case core.EventMIDI_PitchWheel:
			// get the pitch bend value
			m.bend = core.MIDI_PitchBend(me.GetPitchWheel())
			// update all active voices
			for i := range m.voice {
				v := &m.voice[i]
				if v.module != nil {
					f := core.MIDI_ToFrequency(float32(v.note) + m.bend)
					core.SendEventFloatID(v.module, m.freq, f)
				}
			}
		default:
			log.Info.Printf("unhandled midi event %s", me)
		}
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
