//-----------------------------------------------------------------------------
/*

Left/Right Pan and Volume Module

*/
//-----------------------------------------------------------------------------

package audio

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/log"
)

//-----------------------------------------------------------------------------

const (
	pan_port_null = iota
	pan_port_volume
	pan_port_pan
)

// Info returns the module information.
func (m *panModule) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "pan",
		In: []core.PortInfo{
			{"in", "input", core.PortType_AudioBuffer, 0},
			{"volume", "volume (0..1)", core.PortType_EventFloat, pan_port_volume},
			{"pan", "left/right pan (0..1)", core.PortType_EventFloat, pan_port_pan},
		},
		Out: []core.PortInfo{
			{"out_left", "left channel output", core.PortType_AudioBuffer, 0},
			{"out_right", "right channel output", core.PortType_AudioBuffer, 0},
		},
	}
}

//-----------------------------------------------------------------------------

type panModule struct {
	vol          float32 // overall volume
	pan          float32 // pan value 0 == left, 1 == right
	vol_l, vol_r float32 // left right channel volume
}

// NewPan returns a left/right pan and volume module.
func NewPan() core.Module {
	log.Info.Printf("")
	return &panModule{}
}

// Stop and performs any cleanup of a module.
func (m *panModule) Stop() {
	log.Info.Printf("")
}

//-----------------------------------------------------------------------------
// Events

func (m *panModule) set() {
	// Use sin/cos so that l*l + r*r = K (constant power)
	m.vol_l = m.vol * core.Cos(m.pan)
	m.vol_r = m.vol * core.Sin(m.pan)
}

// Event processes a module event.
func (m *panModule) Event(e *core.Event) {
	fe := e.GetEventFloat()
	if fe != nil {
		switch fe.Id {
		case pan_port_volume:
			log.Info.Printf("set volume %f", fe.Val)
			// convert to a linear volume
			m.vol = core.Pow2(core.Clamp(fe.Val, 0, 1)) - 1
			m.set()
		case pan_port_pan:
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
	out_l := buf[1]
	out_r := buf[2]
	// left
	out_l.Copy(in)
	out_l.MulScalar(m.vol_l)
	// right
	out_r.Copy(in)
	out_r.MulScalar(m.vol_r)
}

// Active return true if the module has non-zero output.
func (m *panModule) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
