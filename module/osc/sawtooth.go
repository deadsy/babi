//-----------------------------------------------------------------------------
/*

Sawtooth Oscillator Module

*/
//-----------------------------------------------------------------------------

package osc

import (
	"fmt"

	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

const (
	sawPortNull = iota
	sawPortFrequency
)

// Info returns the module information.
func (m *sawModule) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "saw",
		In: []core.PortInfo{
			{"frequency", "frequency (Hz)", core.PortTypeFloat, sawPortFrequency},
		},
		Out: []core.PortInfo{
			{"out", "output", core.PortTypeAudioBuffer, 0},
		},
	}
}

//-----------------------------------------------------------------------------

type sawType int

const (
	sawTypeNull sawType = iota
	sawTypeBasic
	sawTypeBLEP
)

type sawModule struct {
	synth *core.Synth // top-level synth
	stype sawType     // sawtooth type
	freq  float32     // base frequency
	x     uint32      // phase position
	xstep uint32      // phase step per sample
}

func newSawtooth(s *core.Synth, stype sawType) core.Module {
	return &sawModule{
		synth: s,
		stype: stype,
	}
}

// NewSawtoothBasic returns a non bandwidth limited sawtooth oscillator.
func NewSawtoothBasic(s *core.Synth) core.Module {
	log.Info.Printf("")
	return newSawtooth(s, sawTypeBasic)
}

// NewSawtoothBLEP returns a bandwidth limited sawtooth oscillator.
func NewSawtoothBLEP(s *core.Synth) core.Module {
	log.Info.Printf("")
	return newSawtooth(s, sawTypeBLEP)
}

// Child returns the child modules of this module.
func (m *sawModule) Child() []core.Module {
	return nil
}

// Stop performs any cleanup of a module.
func (m *sawModule) Stop() {
}

//-----------------------------------------------------------------------------
// Events

// Event processes a module event.
func (m *sawModule) Event(e *core.Event) {
	fe := e.GetEventFloat()
	if fe != nil {
		val := fe.Val
		switch fe.ID {
		case sawPortFrequency: // set the oscillator frequency
			log.Info.Printf("set frequency %f", val)
			m.freq = val
			m.xstep = uint32(val * core.FrequencyScale)
		default:
			log.Info.Printf("bad port number %d", fe.ID)
		}
	}
}

//-----------------------------------------------------------------------------

func (m *sawModule) generateBasic(out *core.Buf) {
	for i := 0; i < len(out); i++ {
		out[i] = (2.0/float32(core.FullCycle))*float32(m.x) - 1.0
		// step the phase
		m.x += m.xstep
	}
}

func (m *sawModule) generateBLEP(out *core.Buf) {
	// TODO
}

// Process runs the module DSP.
func (m *sawModule) Process(buf ...*core.Buf) {
	out := buf[0]
	switch m.stype {
	case sawTypeBasic:
		m.generateBasic(out)
	case sawTypeBLEP:
		m.generateBLEP(out)
	default:
		panic(fmt.Sprintf("bad sawtooth type %d", m.stype))
	}
}

// Active returns true if the module has non-zero output.
func (m *sawModule) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
