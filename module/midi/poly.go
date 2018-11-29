//-----------------------------------------------------------------------------
/*

Polyphonic Module

Manage concurrent instances (voices) of a given sub-module.

Note: The single channel output is the sum of outputs from each single channel voice.

*/
//-----------------------------------------------------------------------------

package midi

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

var polyMidiInfo = core.ModuleInfo{
	Name: "polyMidi",
	In: []core.PortInfo{
		{"midi", "midi input", core.PortTypeMIDI, polyMidiIn},
	},
	Out: []core.PortInfo{
		{"out", "output", core.PortTypeAudio, nil},
	},
}

// Info returns the module information.
func (m *polyMidi) Info() *core.ModuleInfo {
	return &m.info
}

//-----------------------------------------------------------------------------

type voiceInfo struct {
	note   uint8       // midi note value
	module core.Module // voice module
}

type polyMidi struct {
	info  core.ModuleInfo                 // module info
	ch    uint8                           // MIDI channel
	sm    func(s *core.Synth) core.Module // new function for voice sub-module
	voice []voiceInfo                     // voices
	idx   int                             // round-robin index for voice slice
	bend  float32                         // pitch bending value (for all voices)
}

// NewPoly returns a MIDI polyphonic voice control module.
func NewPoly(s *core.Synth, ch uint8, sm func(s *core.Synth) core.Module, maxvoices uint) core.Module {
	log.Info.Printf("")
	m := &polyMidi{
		info:  polyMidiInfo,
		ch:    ch,
		sm:    sm,
		voice: make([]voiceInfo, maxvoices),
	}
	return s.Register(m)
}

// Return the child modules.
func (m *polyMidi) Child() []core.Module {
	var children []core.Module
	for i := range m.voice {
		if m.voice[i].module != nil {
			children = append(children, m.voice[i].module)
		}
	}
	return children
}

// Stop performs any cleanup of a module.
func (m *polyMidi) Stop() {
}

//-----------------------------------------------------------------------------
// Events

// voiceLookup returns the voice for this MIDI note (or nil).
func (m *polyMidi) voiceLookup(note uint8) *voiceInfo {
	for i := range m.voice {
		if m.voice[i].module != nil && m.voice[i].note == note {
			return &m.voice[i]
		}
	}
	return nil
}

// voiceAlloc allocates a new subpatch voice for a MIDI note.
func (m *polyMidi) voiceAlloc(note uint8) *voiceInfo {
	log.Info.Printf("note %d", note)
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
	v.module = m.sm(m.Info().Synth)
	// set the voice frequency
	f := core.MIDIToFrequency(float32(v.note) + m.bend)
	core.SendEventFloat(v.module, "frequency", f)
	return v
}

func polyMidiIn(cm core.Module, e *core.Event) {
	m := cm.(*polyMidi)
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
			// perhaps the voices can use this MIDI event...
			for i := range m.voice {
				v := &m.voice[i]
				if v.module != nil {
					core.SendEvent(v.module, "midi", e)
				}
			}
		}
	}
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *polyMidi) Process(buf ...*core.Buf) {
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
func (m *polyMidi) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
