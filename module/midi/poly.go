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
	info    core.ModuleInfo                 // module info
	ch      uint8                           // MIDI channel
	sm      func(s *core.Synth) core.Module // new function for voice sub-module
	voice   []voiceInfo                     // voices
	idx     int                             // round-robin index for voice slice
	bend    float32                         // pitch bending value (for all voices)
	ccCache [128]uint8                      // cache of cc values
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
	// initialise the cc cache to an unused value
	for num := range m.ccCache {
		m.ccCache[num] = 0xff
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
	// send the new voice the cached cc values
	for num := range m.ccCache {
		val := m.ccCache[num]
		if val != 0xff {
			core.EventInMidiCC(v.module, "midi", uint8(num), val)
		}
	}
	// set the voice note
	core.EventInFloat(v.module, "note", float32(v.note)+m.bend)
	return v
}

func polyMidiIn(cm core.Module, e *core.Event) {
	m := cm.(*polyMidi)
	me := e.GetEventMIDIChannel(m.ch)
	if me != nil {
		switch me.GetType() {
		case core.EventMIDINoteOn:
			v := m.voiceLookup(me.GetNote())
			vel := me.GetVelocityFloat()
			if v != nil {
				// note: vel=0 is the same as note off (gate=0).
				core.EventInFloat(v.module, "gate", vel)
			} else {
				if vel != 0 {
					v := m.voiceAlloc(me.GetNote())
					if v != nil {
						core.EventInFloat(v.module, "gate", vel)
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
				core.EventInFloat(v.module, "gate", 0)
			}
		case core.EventMIDIPitchWheel:
			// get the pitch bend value
			m.bend = core.MIDIPitchBend(me.GetPitchWheel())
			// update all active voices
			for i := range m.voice {
				v := &m.voice[i]
				if v.module != nil {
					core.EventInFloat(v.module, "note", float32(v.note)+m.bend)
				}
			}
		case core.EventMIDIControlChange:
			// cache the cc value
			m.ccCache[me.GetCcNum()&0x7f] = me.GetCcInt()
			// and pass the control change though to the voices...
			fallthrough
		default:
			// perhaps the voices can use this MIDI event...
			for i := range m.voice {
				v := &m.voice[i]
				if v.module != nil {
					core.EventIn(v.module, "midi", e)
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
