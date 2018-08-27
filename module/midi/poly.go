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

var ports = []core.PortInfo{
	{"midi", "midi channel input", core.PortType_MIDI, core.PortDirn_In, nil},
}

//-----------------------------------------------------------------------------

type voiceInfo struct {
	note   uint8       // midi note value
	module core.Module // voice module
}

type polyModule struct {
	submodule func() core.Module // new function for voice sub-module
	voice     []voiceInfo        // voices
	idx       int                // round-robin index for voice slice
	bend      float32            // pitch bending value (for all voices)
}

func NewPoly(sm func() core.Module, maxvoices uint) core.Module {
	log.Info.Printf("")
	return &polyModule{
		submodule: sm,
		voice:     make([]voiceInfo, maxvoices),
	}
}

//-----------------------------------------------------------------------------

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
	v.module.Event(core.NewCtrlEvent(core.CtrlEvent_Frequency, f))
	return v
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

// Stop stops and performs any cleanup of a module.
func (m *polyModule) Stop() {
	log.Info.Printf("")
}

// Ports returns the module port information.
func (m *polyModule) Ports() []core.PortInfo {
	return ports
}

// Event processes a module event.
func (m *polyModule) Event(e *core.Event) {
}

//-----------------------------------------------------------------------------
