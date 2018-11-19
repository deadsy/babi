//-----------------------------------------------------------------------------
/*

Left/Right Pan and Volume Module

Takes a single audio buffer stream as input and outputs left and right channels.
The "pan" control pans the signal between the left and right channels with
constant power. The "volume" control sets the overall power output.

*/
//-----------------------------------------------------------------------------

package mix

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

// Info returns the module information.
func (m *panMix) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "panMix",
		In: []core.PortInfo{
			{"in", "input", core.PortTypeAudioBuffer, nil},
			{"volume", "volume (0..1)", core.PortTypeFloat, panMixVolume},
			{"pan", "left/right pan (0..1)", core.PortTypeFloat, panMixPan},
		},
		Out: []core.PortInfo{
			{"out_left", "left channel output", core.PortTypeAudioBuffer, nil},
			{"out_right", "right channel output", core.PortTypeAudioBuffer, nil},
		},
	}
}

//-----------------------------------------------------------------------------

type panMix struct {
	synth *core.Synth // top-level synth
	vol   float32     // overall volume
	pan   float32     // pan value 0 == left, 1 == right
	volL  float32     // left channel volume
	volR  float32     // right channel volume
}

// NewPan returns a left/right pan and volume module.
func NewPan(s *core.Synth) core.Module {
	log.Info.Printf("")
	return &panMix{
		synth: s,
	}
}

// Return the child modules.
func (m *panMix) Child() []core.Module {
	return nil
}

// Stop and performs any cleanup of a module.
func (m *panMix) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

func (m *panMix) set() {
	// Use sin/cos so that l*l + r*r = K (constant power)
	m.volL = m.vol * core.Cos(m.pan)
	m.volR = m.vol * core.Sin(m.pan)
}

func panMixVolume(cm core.Module, e *core.Event) {
	m := cm.(*panMix)
	vol := core.Clamp(e.GetEventFloat().Val, 0, 1)
	log.Info.Printf("set volume %f", vol)
	// convert to a linear volume
	m.vol = core.Pow2(vol) - 1.0
	m.set()
}

func panMixPan(cm core.Module, e *core.Event) {
	m := cm.(*panMix)
	pan := core.Clamp(e.GetEventFloat().Val, 0, 1)
	log.Info.Printf("set pan %f", pan)
	m.pan = pan * (core.Pi / 2.0)
	m.set()
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *panMix) Process(buf ...*core.Buf) {
	in := buf[0]
	outL := buf[1]
	outR := buf[2]
	// left
	outL.Copy(in)
	outL.MulScalar(m.volL)
	// right
	outR.Copy(in)
	outR.MulScalar(m.volR)
}

// Active return true if the module has non-zero output.
func (m *panMix) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
