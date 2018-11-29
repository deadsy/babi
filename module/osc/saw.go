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

var sawOscInfo = core.ModuleInfo{
	Name: "sawOsc",
	In: []core.PortInfo{
		{"frequency", "frequency (Hz)", core.PortTypeFloat, sawPortFrequency},
	},
	Out: []core.PortInfo{
		{"out", "output", core.PortTypeAudio, nil},
	},
}

// Info returns the module information.
func (m *sawOsc) Info() *core.ModuleInfo {
	return &m.info
}

//-----------------------------------------------------------------------------

type sawType int

const (
	sawTypeNull sawType = iota
	sawTypeBasic
	sawTypeBLEP
)

type sawOsc struct {
	info  core.ModuleInfo // module info
	stype sawType         // sawtooth type
	freq  float32         // base frequency
	x     uint32          // phase position
	xstep uint32          // phase step per sample
}

func newSawtooth(s *core.Synth, stype sawType) core.Module {
	m := &sawOsc{
		info:  sawOscInfo,
		stype: stype,
	}
	return s.Register(m)
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
func (m *sawOsc) Child() []core.Module {
	return nil
}

// Stop performs any cleanup of a module.
func (m *sawOsc) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

func sawPortFrequency(cm core.Module, e *core.Event) {
	m := cm.(*sawOsc)
	frequency := core.ClampLo(e.GetEventFloat().Val, 0)
	log.Info.Printf("set frequency %f Hz", frequency)
	m.freq = frequency
	m.xstep = uint32(frequency * core.FrequencyScale)
}

//-----------------------------------------------------------------------------

func (m *sawOsc) generateBasic(out *core.Buf) {
	for i := 0; i < len(out); i++ {
		out[i] = (2.0/float32(core.FullCycle))*float32(m.x) - 1.0
		// step the phase
		m.x += m.xstep
	}
}

func (m *sawOsc) generateBLEP(out *core.Buf) {
	// TODO
}

// Process runs the module DSP.
func (m *sawOsc) Process(buf ...*core.Buf) {
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
func (m *sawOsc) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
