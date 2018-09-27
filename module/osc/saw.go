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

// Info returns the module information.
func (m *sawModule) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "saw",
		In: []core.PortInfo{
			{"frequency", "frequency (Hz)", core.PortTypeFloat, sawPortFrequency},
		},
		Out: []core.PortInfo{
			{"out", "output", core.PortTypeAudioBuffer, nil},
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
// Port Events

func sawPortFrequency(cm core.Module, e *core.Event) {
	m := cm.(*sawModule)
	frequency := core.ClampLo(e.GetEventFloat().Val, 0)
	log.Info.Printf("set frequency %f Hz", frequency)
	m.freq = frequency
	m.xstep = uint32(frequency * core.FrequencyScale)
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
