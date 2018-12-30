//-----------------------------------------------------------------------------
/*

Metronome, MIDI out, sequencer testing.

*/
//-----------------------------------------------------------------------------

package main

import (
	"os"
	"os/signal"

	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/midi"
	"github.com/deadsy/babi/module/seq"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

const ch = 1

var metronome = []seq.Op{
	seq.OpNote(ch, 69, 100, 4),
	seq.OpRest(12),
	seq.OpNote(ch, 60, 100, 4),
	seq.OpRest(12),
	seq.OpNote(ch, 60, 100, 4),
	seq.OpRest(12),
	seq.OpNote(ch, 60, 100, 4),
	seq.OpRest(12),
	seq.OpLoop(),
}

//-----------------------------------------------------------------------------
// Top level Metronome Patch

var metroInfo = core.ModuleInfo{
	Name: "metro",
	In: []core.PortInfo{
		{"midi", "midi input", core.PortTypeMIDI, metroMidiIn},
	},
	Out: []core.PortInfo{
		{"midi", "midi output", core.PortTypeMIDI, nil},
	},
}

// Info returns the general module information.
func (m *metro) Info() *core.ModuleInfo {
	return &m.info
}

type metro struct {
	info core.ModuleInfo // module info
	seq  core.Module     // sequencer
}

// NewMetro returns a metronome module.
func NewMetro(s *core.Synth, prog []seq.Op) core.Module {
	log.Info.Printf("")

	sx := seq.NewSequencer(s, prog)

	// defaults
	core.EventInFloat(sx, "bpm", 120.0)
	core.EventInInt(sx, "ctrl", seq.CtrlStart)

	// monitor the MIDI events
	mon := midi.NewMonitor(s, ch)
	core.Connect(sx, "midi", mon, "midi")

	m := &metro{
		info: metroInfo,
		seq:  sx,
	}
	return s.Register(m)
}

// Child returns the child modules of this module.
func (m *metro) Child() []core.Module {
	return []core.Module{m.seq}
}

// Stop performs any cleanup of a module.
func (m *metro) Stop() {
}

// Port Events

func metroMidiIn(cm core.Module, e *core.Event) {
	m := cm.(*metro)
	_ = m
	// TODO process a CC for bpm control, etc.
}

// Process runs the module DSP. Return true for non-zero output.
func (m *metro) Process(buf ...*core.Buf) bool {
	m.seq.Process(nil)
	return false
}

//-----------------------------------------------------------------------------

func main() {
	s := core.NewSynth()

	// create the metronome patch
	p := NewMetro(s, metronome)

	// set the root patch for the synth
	s.SetPatch(p)

	// start the jack client
	err := s.StartJack("metro")
	if err != nil {
		log.Error.Printf("%s", err)
		s.Close()
		os.Exit(1)
	}

	// signal handling
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	<-signals

	s.Close()
	os.Exit(0)
}

//-----------------------------------------------------------------------------
