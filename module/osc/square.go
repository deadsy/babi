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

const (
	sqrPortNull = iota
	sqrPortFrequency
	sqrPortDuty
)

// Info returns the module information.
func (m *sqrModule) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "square",
		In: []core.PortInfo{
			{"frequency", "frequency (Hz)", core.PortTypeFloat, sqrPortFrequency},
			{"duty", "duty cycle (0..1)", core.PortTypeFloat, sqrPortDuty},
		},
		Out: []core.PortInfo{
			{"out", "output", core.PortTypeAudioBuffer, 0},
		},
	}
}

//-----------------------------------------------------------------------------

const sqrFullCycle = float32(1 << 32)

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

// NewX returns an X module.
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
// Events

// Event processes a module event.
func (m *sqrModule) Event(e *core.Event) {
	fe := e.GetEventFloat()
	if fe != nil {
		val := fe.Val
		switch fe.ID {
		case sqrPortDuty: // set the duty cycle
			log.Info.Printf("set duty cycle %f", val)
			duty := core.Clamp(val, 0, 1)
			m.tp = uint32(sqrFullCycle * core.Map(duty, 0.05, 0.5))
		case sqrPortFrequency: // set the oscillator frequency
			log.Info.Printf("set frequency %f", val)
			m.freq = val
			m.xstep = uint32(val * core.FrequencyScale)
		default:
			log.Info.Printf("bad port number %d", fe.ID)
		}
	}
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
