//-----------------------------------------------------------------------------
/*

Use a mode selector to multiplex a limited number of control channels to a greater
number of control events.

*/
//-----------------------------------------------------------------------------

package midi

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

var modalMidiInfo = core.ModuleInfo{
	Name: "modalMidi",
	In: []core.PortInfo{
		{"midi", "midi input", core.PortTypeMIDI, modalMidiIn},
		{"mode", "mode selector", core.PortTypeInt, modalMidiMode},
	},
	Out: nil,
}

// Info returns the general module information.
func (m *modalMidi) Info() *core.ModuleInfo {
	return &m.info
}

//-----------------------------------------------------------------------------

type modalMidi struct {
	info core.ModuleInfo // module info
}

// NewModal returns an modela midi control selector.
func NewModal(s *core.Synth, ch, cc uint8, nControls, nModes int) core.Module {
	log.Info.Printf("midi ch %d cc %d %dx%d controls", ch, cc, nModes, nControls)
	m := &modalMidi{
		info: modalMidiInfo,
	}
	return s.Register(m)
}

// Child returns the child modules of this module.
func (m *modalMidi) Child() []core.Module {
	return nil
}

// Stop performs any cleanup of a module.
func (m *modalMidi) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

func modalMidiIn(cm core.Module, e *core.Event) {
}

func modalMidiMode(cm core.Module, e *core.Event) {
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *modalMidi) Process(buf ...*core.Buf) {
}

// Active returns true if the module has non-zero output.
func (m *modalMidi) Active() bool {
	return false
}

//-----------------------------------------------------------------------------
