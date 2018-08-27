//-----------------------------------------------------------------------------
/*

Noise Generator Module

https://noisehack.com/generate-noise-web-audio-api/
http://www.musicdsp.org/files/pink.txt
https://en.wikipedia.org/wiki/Pink_noise
https://en.wikipedia.org/wiki/White_noise
https://en.wikipedia.org/wiki/Brownian_noise

*/
//-----------------------------------------------------------------------------

package noise

import (
	"fmt"

	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/log"
)

//-----------------------------------------------------------------------------

type noiseType int

const (
	noiseType_null noiseType = iota
	noiseType_pink1
	noiseType_pink2
	noiseType_white
	noiseType_brown
)

type noiseModule struct {
	ntype                      noiseType  // noise type
	r                          *core.Rand // random state
	b0, b1, b2, b3, b4, b5, b6 float32    // state variables
}

func newNoise(ntype noiseType) core.Module {
	log.Info.Printf("")
	return &noiseModule{
		ntype: ntype,
	}
}

// NewWhite returns a white noise generator module.
func NewWhite() core.Module {
	return newNoise(noiseType_white)
}

// NewBrown returns a brown noise generator module.
func NewBrown() core.Module {
	return newNoise(noiseType_brown)
}

// NewPink1 returns a pink noise generator module.
func NewPink1() core.Module {
	return newNoise(noiseType_pink1)
}

// NewPink2 returns a pink noise generator module.
func NewPink2() core.Module {
	return newNoise(noiseType_pink2)
}

// Stop and performs any cleanup of a module.
func (m *noiseModule) Stop() {
	log.Info.Printf("")
}

//-----------------------------------------------------------------------------
// Ports

var noisePorts = []core.PortInfo{
	{"out", "output", core.PortType_Buf, core.PortDirn_Out, nil},
}

// Ports returns the module port information.
func (m *noiseModule) Ports() []core.PortInfo {
	return noisePorts
}

//-----------------------------------------------------------------------------
// Events

// Event processes a module event.
func (m *noiseModule) Event(e *core.Event) {
	// do nothing
}

//-----------------------------------------------------------------------------

func (m *noiseModule) generate_white(out *core.Buf) {
}

func (m *noiseModule) generate_brown(out *core.Buf) {
}

func (m *noiseModule) generate_pink1(out *core.Buf) {
}

func (m *noiseModule) generate_pink2(out *core.Buf) {
}

// Process runs the module DSP.
func (m *noiseModule) Process(buf []*core.Buf) {
	out := buf[0]
	switch m.ntype {
	case noiseType_white:
		m.generate_white(out)
	case noiseType_pink1:
		m.generate_pink1(out)
	case noiseType_pink2:
		m.generate_pink2(out)
	case noiseType_brown:
		m.generate_brown(out)
	default:
		panic(fmt.Sprintf("bad noise type %d", m.ntype))
	}
}

// Active return true if the module has non-zero output.
func (m *noiseModule) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
