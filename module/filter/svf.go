//-----------------------------------------------------------------------------
/*

State Variable Filters

*/
//-----------------------------------------------------------------------------

package filter

import (
	"fmt"

	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

var svFilterInfo = core.ModuleInfo{
	Name: "svFilter",
	In: []core.PortInfo{
		{"in", "input", core.PortTypeAudio, nil},
		{"cutoff", "cutoff frequency (Hz)", core.PortTypeFloat, svfPortCutoff},
		{"resonance", "resonance (0..1)", core.PortTypeFloat, svfPortResonance},
	},
	Out: []core.PortInfo{
		{"out", "output", core.PortTypeAudio, nil},
	},
}

// Info returns the module information.
func (m *svFilter) Info() *core.ModuleInfo {
	return &m.info
}

//-----------------------------------------------------------------------------

type svfType int

const (
	svfTypeNull svfType = iota
	svfTypeHC
	svfTypeTrapezoidal
)

type svFilter struct {
	info  core.ModuleInfo // module info
	ftype svfType         // filter type
	// svfTypeHC
	kf float32 // constant for cutoff frequency
	kq float32 // constant for filter resonance
	bp float32 // bandpass state variable
	lp float32 // low pass state variable
	// svfTypeTrapezoidal
	g     float32 // constant for cutoff frequency
	k     float32 // constant for filter resonance
	ic1eq float32 // state variable
	ic2eq float32 // state variable
}

func newSVF(s *core.Synth, t svfType) core.Module {
	m := &svFilter{
		info:  svFilterInfo,
		ftype: t,
	}
	return s.Register(m)
}

// NewSVFilterHC returns a state variable filter.
// See: Hal Chamberlin's "Musical Applications of Microprocessors" pp.489-492.
func NewSVFilterHC(s *core.Synth) core.Module {
	log.Info.Printf("")
	return newSVF(s, svfTypeHC)
}

// NewSVFilterTrapezoidal returns a state variable filter.
// See: https://cytomic.com/files/dsp/SvfLinearTrapOptimised2.pdf
func NewSVFilterTrapezoidal(s *core.Synth) core.Module {
	log.Info.Printf("")
	return newSVF(s, svfTypeTrapezoidal)
}

// Child returns the child modules of this module.
func (m *svFilter) Child() []core.Module {
	return nil
}

// Stop performs any cleanup of a module.
func (m *svFilter) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

func svfPortCutoff(cm core.Module, e *core.Event) {
	m := cm.(*svFilter)
	cutoff := core.Clamp(e.GetEventFloat().Val, 0, 0.5*core.AudioSampleFrequency)
	log.Info.Printf("set cutoff frequency %f Hz", cutoff)
	switch m.ftype {
	case svfTypeHC:
		m.kf = 2.0 * core.Sin(core.Pi*cutoff*core.AudioSamplePeriod)
	case svfTypeTrapezoidal:
		m.g = core.Tan(core.Pi * cutoff * core.AudioSamplePeriod)
	default:
		panic(fmt.Sprintf("bad filter type %d", m.ftype))
	}
}

func svfPortResonance(cm core.Module, e *core.Event) {
	m := cm.(*svFilter)
	resonance := core.Clamp(e.GetEventFloat().Val, 0, 1)
	log.Info.Printf("set resonance %f", resonance)
	switch m.ftype {
	case svfTypeHC:
		m.kq = 2.0 - 2.0*resonance
	case svfTypeTrapezoidal:
		m.k = 2.0 - 2.0*resonance
	default:
		panic(fmt.Sprintf("bad filter type %d", m.ftype))
	}
}

//-----------------------------------------------------------------------------

func (m *svFilter) filterHC(in, out *core.Buf) {
	lp := m.lp
	bp := m.bp
	kf := m.kf
	kq := m.kq
	for i := 0; i < len(out); i++ {
		lp += kf * bp
		hp := in[i] - lp - (kq * bp)
		bp += kf * hp
		out[i] = lp
	}
	// update the state variables
	m.lp = lp
	m.bp = bp
}

func (m *svFilter) filterTrapezoidal(in, out *core.Buf) {
	ic1eq := m.ic1eq
	ic2eq := m.ic2eq
	a1 := 1.0 / (1.0 + (m.g * (m.g + m.k)))
	a2 := m.g * a1
	a3 := m.g * a2
	for i := 0; i < len(out); i++ {
		v0 := in[i]
		v3 := v0 - ic2eq
		v1 := (a1 * ic1eq) + (a2 * v3)
		v2 := ic2eq + (a2 * ic1eq) + (a3 * v3)
		ic1eq = (2.0 * v1) - ic1eq
		ic2eq = (2.0 * v2) - ic2eq
		out[i] = v2 // low
		//low := v2
		//band := v1
		//high := v0 - (m.k * v1) - v2
		//notch := v0 - (m.k * v1)
		//peak := v0 - (m.k * v1) - (2.0 * v2)
		//all := v0 - (2.0 * m.k * v1)
	}
	// update the state variables
	m.ic1eq = ic1eq
	m.ic2eq = ic2eq
}

// Process runs the module DSP.
func (m *svFilter) Process(buf ...*core.Buf) {
	in := buf[0]
	out := buf[1]
	switch m.ftype {
	case svfTypeHC:
		m.filterHC(in, out)
	case svfTypeTrapezoidal:
		m.filterTrapezoidal(in, out)
	default:
		panic(fmt.Sprintf("bad filter type %d", m.ftype))
	}
}

// Active returns true if the module has non-zero output.
func (m *svFilter) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
