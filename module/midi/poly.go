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
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

// Info returns the module information.
func (m *polyModule) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "poly",
		In: []core.PortInfo{
			{"midi_in", "midi input", core.PortTypeMIDI, polyPortMidiIn},
		},
		Out: []core.PortInfo{
			{"out", "output", core.PortTypeAudioBuffer, nil},
		},
	}
}

//-----------------------------------------------------------------------------

type voiceInfo struct {
	note   uint8       // midi note value
	module core.Module // voice module
}

type polyModule struct {
	synth     *core.Synth        // top-level synth
	ch        uint8              // MIDI channel
	submodule func() core.Module // new function for voice sub-module
	voice     []voiceInfo        // voices
	idx       int                // round-robin index for voice slice
	bend      float32            // pitch bending value (for all voices)
}

// NewPoly returns a MIDI polyphonic voice control module.
func NewPoly(s *core.Synth, ch uint8, sm func() core.Module, maxvoices uint) core.Module {
	log.Info.Printf("")
	return &polyModule{
		synth:     s,
		ch:        ch,
		submodule: sm,
		voice:     make([]voiceInfo, maxvoices),
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
	m.idx++
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
	f := core.MIDIToFrequency(float32(v.note) + m.bend)
	core.SendEventFloat(v.module, "frequency", f)
	return v
}

func polyPortMidiIn(cm core.Module, e *core.Event) {
	m := cm.(*polyModule)
	me := e.GetEventMIDIChannel(m.ch)
	if me != nil {
		switch me.GetType() {
		case core.EventMIDINoteOn:
			v := m.voiceLookup(me.GetNote())
			vel := core.MIDIMap(me.GetVelocity(), 0, 1)
			if v != nil {
				// note: vel=0 is the same as note off (gate=0).
				core.SendEventFloat(v.module, "gate", vel)
			} else {
				if vel != 0 {
					v := m.voiceAlloc(me.GetNote())
					if v != nil {
						core.SendEventFloat(v.module, "gate", vel)
					} else {
						log.Info.Printf("unable to allocate new voice")
					}
				}
			}
		case core.EventMIDINoteOff:
			v := m.voiceLookup(me.GetNote())
			if v != nil {
				// send a note off control event
				// ignoring the note off velocity (for now)
				core.SendEventFloat(v.module, "gate", 0)
			}
		case core.EventMIDIPitchWheel:
			// get the pitch bend value
			m.bend = core.MIDIPitchBend(me.GetPitchWheel())
			// update all active voices
			for i := range m.voice {
				v := &m.voice[i]
				if v.module != nil {
					f := core.MIDIToFrequency(float32(v.note) + m.bend)
					core.SendEventFloat(v.module, "frequency", f)
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
