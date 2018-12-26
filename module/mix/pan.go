//-----------------------------------------------------------------------------
/*

Left/Right Pan and Volume Module

Takes a single audio buffer stream as input and outputs left and right channels.

*/
//-----------------------------------------------------------------------------

package mix

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

var panMixInfo = core.ModuleInfo{
	Name: "panMix",
	In: []core.PortInfo{
		{"in", "input", core.PortTypeAudio, nil},
		{"midi", "midi input", core.PortTypeMIDI, panMixMidiIn},
		{"vol", "volume (0..1)", core.PortTypeFloat, panMixVolume},
		{"pan", "left/right pan (0..1)", core.PortTypeFloat, panMixPan},
	},
	Out: []core.PortInfo{
		{"out0", "left channel output", core.PortTypeAudio, nil},
		{"out1", "right channel output", core.PortTypeAudio, nil},
	},
}

// Info returns the module information.
func (m *panMix) Info() *core.ModuleInfo {
	return &m.info
}

//-----------------------------------------------------------------------------

type panMix struct {
	info  core.ModuleInfo // module info
	ch    uint8           // MIDI channel
	ccPan uint8           // MIDI CC number for pan control
	ccVol uint8           // MIDI CC number for volume control
	vol   float32         // overall volume
	pan   float32         // pan value 0 == left, 1 == right
	volL  float32         // left channel volume
	volR  float32         // right channel volume
}

// NewPan returns a left/right pan and volume module.
func NewPan(s *core.Synth, ch, cc uint8) core.Module {
	log.Info.Printf("")
	m := &panMix{
		info:  panMixInfo,
		ch:    ch,
		ccPan: cc,
		ccVol: cc + 1,
	}
	return s.Register(m)
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

func (m *panMix) setVol(vol float32) {
	log.Info.Printf("set volume %f", vol)
	// convert to a linear volume
	m.vol = core.Pow2(vol) - 1.0
	m.set()
}

func (m *panMix) setPan(pan float32) {
	log.Info.Printf("set pan %f", pan)
	m.pan = pan * (core.Pi / 2.0)
	m.set()
}

func panMixVolume(cm core.Module, e *core.Event) {
	m := cm.(*panMix)
	m.setVol(core.Clamp(e.GetEventFloat().Val, 0, 1))
}

func panMixPan(cm core.Module, e *core.Event) {
	m := cm.(*panMix)
	m.setPan(core.Clamp(e.GetEventFloat().Val, 0, 1))
}

func panMixMidiIn(cm core.Module, e *core.Event) {
	m := cm.(*panMix)
	me := e.GetEventMIDIChannel(m.ch)
	if me != nil {
		if me.GetType() == core.EventMIDIControlChange {
			switch me.GetCcNum() {
			case m.ccVol:
				m.setVol(me.GetCcFloat())
			case m.ccPan:
				m.setPan(me.GetCcFloat())
			default:
				// ignore
			}
		}
	}
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *panMix) Process(buf ...*core.Buf) bool {
	in := buf[0]
	out0 := buf[1]
	out1 := buf[2]
	// left
	out0.Copy(in)
	out0.MulScalar(m.volL)
	// right
	out1.Copy(in)
	out1.MulScalar(m.volR)
	return true
}

//-----------------------------------------------------------------------------
