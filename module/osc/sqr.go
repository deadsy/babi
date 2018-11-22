//-----------------------------------------------------------------------------
/*

Square Wave Oscillator Module

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
func (m *sqrModule) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "sqr",
		In: []core.PortInfo{
			{"frequency", "frequency (Hz)", core.PortTypeFloat, sqrPortFrequency},
			{"duty", "duty cycle (0..1)", core.PortTypeFloat, sqrPortDuty},
		},
		Out: []core.PortInfo{
			{"out", "output", core.PortTypeAudio, nil},
		},
	}
}

//-----------------------------------------------------------------------------

type sqrType int

const (
	sqrTypeNull sqrType = iota
	sqrTypeBasic
	sqrTypeBLEP
)

type sqrModule struct {
	synth *core.Synth // top-level synth
	stype sqrType     // square type
	tp    uint32      // 1/0 transition point
	freq  float32     // base frequency
	x     uint32      // phase position
	xstep uint32      // phase step per sample
}

func newSquare(s *core.Synth, stype sqrType) core.Module {
	return &sqrModule{
		synth: s,
		stype: stype,
	}
}

// NewSquareBasic returns a non bandwidth limited square wave oscillator.
func NewSquareBasic(s *core.Synth) core.Module {
	log.Info.Printf("")
	return newSquare(s, sqrTypeBasic)
}

// NewSquareBLEP returns a bandwidth limited square wave oscillator.
func NewSquareBLEP(s *core.Synth) core.Module {
	log.Info.Printf("")
	return newSquare(s, sqrTypeBLEP)
}

// Child returns the child modules of this module.
func (m *sqrModule) Child() []core.Module {
	return nil
}

// Stop performs any cleanup of a module.
func (m *sqrModule) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

func sqrPortFrequency(cm core.Module, e *core.Event) {
	m := cm.(*sqrModule)
	frequency := core.ClampLo(e.GetEventFloat().Val, 0)
	log.Info.Printf("set frequency %f Hz", frequency)
	m.freq = frequency
	m.xstep = uint32(frequency * core.FrequencyScale)
}

func sqrPortDuty(cm core.Module, e *core.Event) {
	m := cm.(*sqrModule)
	duty := core.Clamp(e.GetEventFloat().Val, 0, 1)
	log.Info.Printf("set duty cycle %f", duty)
	m.tp = uint32(float32(core.FullCycle) * core.Map(duty, 0.05, 0.5))
}

//-----------------------------------------------------------------------------

func (m *sqrModule) generateBasic(out *core.Buf) {
	for i := 0; i < len(out); i++ {
		// what portion of the cycle are we in?
		if m.x < m.tp {
			out[i] = 1
		} else {
			out[i] = -1
		}
		// step the phase
		m.x += m.xstep
	}
}

func (m *sqrModule) generateBLEP(out *core.Buf) {
	// TODO
}

// Process runs the module DSP.
func (m *sqrModule) Process(buf ...*core.Buf) {
	out := buf[0]
	switch m.stype {
	case sqrTypeBasic:
		m.generateBasic(out)
	case sqrTypeBLEP:
		m.generateBLEP(out)
	default:
		panic(fmt.Sprintf("bad square type %d", m.stype))
	}
}

// Active returns true if the module has non-zero output.
func (m *sqrModule) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
