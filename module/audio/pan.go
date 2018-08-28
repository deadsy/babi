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
// Ports

var panPorts = []core.PortInfo{
	{"in", "input", core.PortType_AudioBuffer, core.PortDirn_In},
	{"out_l", "left channel output", core.PortType_AudioBuffer, core.PortDirn_Out},
	{"out_r", "right channel output", core.PortType_AudioBuffer, core.PortDirn_Out},
	{"vol", "volume (0..1)", core.PortType_EventFloat32, core.PortDirn_In},
	{"pan", "left/right pan (0..1)", core.PortType_EventFloat32, core.PortDirn_In},
}

// Ports returns the module port information.
func (m *panModule) Ports() []core.PortInfo {
	return panPorts
}

//-----------------------------------------------------------------------------
// Events

/*
func (p *Pan) set() {
	p.vol_l = p.vol * core.Cos(p.pan)
	p.vol_r = p.vol * core.Sin(p.pan)
}

func (p *Pan) SetVol(vol float32) {
	// convert to a linear volume
	p.vol = core.Pow2(core.Clamp(vol, 0, 1)) - 1
	p.set()
}

func (p *Pan) SetPan(pan float32) {
	// Use sin/cos so that l*l + r*r = K (constant power)
	p.pan = core.Clamp(pan, 0, 1) * (core.PI / 2)
	p.set()
}
*/

// Event processes a module event.
func (m *panModule) Event(e *core.Event) {
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
