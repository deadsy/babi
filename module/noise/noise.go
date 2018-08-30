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

// Info returns the module information.
func (m *noiseModule) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "noise",
		In:   nil,
		Out: []core.PortInfo{
			{"out", "output", core.PortType_AudioBuffer, 0},
		},
	}
}

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
	return &noiseModule{
		ntype: ntype,
	}
}

// NewWhite returns a white noise generator module.
// white noise (spectral density = k)
func NewWhite() core.Module {
	log.Info.Printf("")
	return newNoise(noiseType_white)
}

// NewBrown returns a brown noise generator module.
// brown noise (spectral density = k/f*f)
func NewBrown() core.Module {
	log.Info.Printf("")
	return newNoise(noiseType_brown)
}

// NewPink1 returns a pink noise generator module.
// pink noise (spectral density = k/f): fast, inaccurate version
func NewPink1() core.Module {
	log.Info.Printf("")
	return newNoise(noiseType_pink1)
}

// NewPink2 returns a pink noise generator module.
// pink noise (spectral density = k/f): slow, accurate version
func NewPink2() core.Module {
	log.Info.Printf("")
	return newNoise(noiseType_pink2)
}

// Stop and performs any cleanup of a module.
func (m *noiseModule) Stop() {
	log.Info.Printf("")
}

//-----------------------------------------------------------------------------
// Events

// Event processes a module event.
func (m *noiseModule) Event(e *core.Event) {
	// do nothing
}

//-----------------------------------------------------------------------------

func (m *noiseModule) generate_white(out *core.Buf) {
	for i := 0; i < len(out); i++ {
		out[i] = m.r.Float()
	}
}

func (m *noiseModule) generate_brown(out *core.Buf) {
	b0 := m.b0
	for i := 0; i < len(out); i++ {
		white := m.r.Float()
		b0 = (b0 + (0.02 * white)) * (1.0 / 1.02)
		out[i] = b0 * (1.0 / 0.38)
	}
	m.b0 = b0
}

func (m *noiseModule) generate_pink1(out *core.Buf) {
	b0 := m.b0
	b1 := m.b1
	b2 := m.b2
	for i := 0; i < len(out); i++ {
		white := m.r.Float()
		b0 = 0.99765*b0 + white*0.0990460
		b1 = 0.96300*b1 + white*0.2965164
		b2 = 0.57000*b2 + white*1.0526913
		pink := b0 + b1 + b2 + white*0.1848
		out[i] = pink * (1.0 / 10.4)
	}
	m.b0 = b0
	m.b1 = b1
	m.b2 = b2
}

func (m *noiseModule) generate_pink2(out *core.Buf) {
	b0 := m.b0
	b1 := m.b1
	b2 := m.b2
	b3 := m.b3
	b4 := m.b4
	b5 := m.b5
	b6 := m.b6
	for i := 0; i < len(out); i++ {
		white := m.r.Float()
		b0 = 0.99886*b0 + white*0.0555179
		b1 = 0.99332*b1 + white*0.0750759
		b2 = 0.96900*b2 + white*0.1538520
		b3 = 0.86650*b3 + white*0.3104856
		b4 = 0.55000*b4 + white*0.5329522
		b5 = -0.7616*b5 - white*0.0168980
		pink := b0 + b1 + b2 + b3 + b4 + b5 + b6 + white*0.5362
		b6 = white * 0.115926
		out[i] = pink * (1.0 / 10.2)
	}
	m.b0 = b0
	m.b1 = b1
	m.b2 = b2
	m.b3 = b3
	m.b4 = b4
	m.b5 = b5
	m.b6 = b6
}

// Process runs the module DSP.
func (m *noiseModule) Process(buf ...*core.Buf) {
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
