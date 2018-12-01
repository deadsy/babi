//-----------------------------------------------------------------------------
/*

MIDI event monitor

Logs the MIDI events on a MIDI channel.

*/
//-----------------------------------------------------------------------------

package midi

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

var monitorMidiInfo = core.ModuleInfo{
	Name: "monitorMidi",
	In: []core.PortInfo{
		{"midi", "midi input", core.PortTypeMIDI, monitorMidiIn},
	},
	Out: nil,
}

// Info returns the general module information.
func (m *monitorMidi) Info() *core.ModuleInfo {
	return &m.info
}

//-----------------------------------------------------------------------------

type monitorMidi struct {
	info core.ModuleInfo // module info
	ch   uint8           // MIDI channel
}

// NewMonitor returns a MIDI monitor module.
func NewMonitor(s *core.Synth, ch uint8) core.Module {
	log.Info.Printf("")
	m := &monitorMidi{
		info: monitorMidiInfo,
		ch:   ch,
	}
	return s.Register(m)
}

// Child returns the child modules of this module.
func (m *monitorMidi) Child() []core.Module {
	return nil
}

// Stop performs any cleanup of a module.
func (m *monitorMidi) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

func monitorMidiIn(cm core.Module, e *core.Event) {
	m := cm.(*monitorMidi)
	me := e.GetEventMIDIChannel(m.ch)
	if me != nil {
		log.Info.Printf("%s", me.String())
	}
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *monitorMidi) Process(buf ...*core.Buf) {
}

// Active returns true if the module has non-zero output.
func (m *monitorMidi) Active() bool {
	return false
}

//-----------------------------------------------------------------------------
