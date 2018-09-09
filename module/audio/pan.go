//-----------------------------------------------------------------------------
/*

Left/Right Pan and Volume Module

Takes a single audio buffer stream as input and outputs left and right channels.
The "pan" control pans the signal between the left and right channels with
constant power. The "volume" control sets the overall power output.

*/
//-----------------------------------------------------------------------------

package audio

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

const (
	panPortNull = iota
	panPortVolume
	panPortPan
)

// Info returns the module information.
func (m *panModule) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "pan",
		In: []core.PortInfo{
			{"in", "input", core.PortType_AudioBuffer, 0},
			{"volume", "volume (0..1)", core.PortType_EventFloat, panPortVolume},
			{"pan", "left/right pan (0..1)", core.PortType_EventFloat, panPortPan},
		},
		Out: []core.PortInfo{
			{"out_left", "left channel output", core.PortType_AudioBuffer, 0},
			{"out_right", "right channel output", core.PortType_AudioBuffer, 0},
		},
	}
}

//-----------------------------------------------------------------------------

type panModule struct {
	vol  float32 // overall volume
	pan  float32 // pan value 0 == left, 1 == right
	volL float32 // left channel volume
	volR float32 // right channel volume
}

// NewPan returns a left/right pan and volume module.
func NewPan() core.Module {
	log.Info.Printf("")
	return &panModule{}
}

// Return the child modules.
func (m *panModule) Child() []core.Module {
	return nil
}

// Stop and performs any cleanup of a module.
func (m *panModule) Stop() {
	log.Info.Printf("")
}

//-----------------------------------------------------------------------------
// Events

func (m *panModule) set() {
	// Use sin/cos so that l*l + r*r = K (constant power)
	m.volL = m.vol * core.Cos(m.pan)
	m.volR = m.vol * core.Sin(m.pan)
}

// Event processes a module event.
func (m *panModule) Event(e *core.Event) {
	fe := e.GetEventFloat()
	if fe != nil {
		switch fe.Id {
		case panPortVolume:
			log.Info.Printf("set volume %f", fe.Val)
			// convert to a linear volume
			m.vol = core.Pow2(core.Clamp(fe.Val, 0, 1)) - 1
			m.set()
		case panPortPan:
			log.Info.Printf("set pan %f", fe.Val)
			m.pan = core.Clamp(fe.Val, 0, 1) * (core.PI / 2)
			m.set()
		default:
			log.Info.Printf("bad port number %d", fe.Id)
		}
	}
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *panModule) Process(buf ...*core.Buf) {
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
func (m *panModule) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
