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

package osc

import (
	"fmt"

	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

var noiseOscInfo = core.ModuleInfo{
	Name: "noiseOsc",
	In:   nil,
	Out: []core.PortInfo{
		{"out", "output", core.PortTypeAudio, nil},
	},
}

// Info returns the module information.
func (m *noiseOsc) Info() *core.ModuleInfo {
	return &m.info
}

//-----------------------------------------------------------------------------

type noiseType int

const (
	noiseTypeNull noiseType = iota
	noiseTypePink1
	noiseTypePink2
	noiseTypeWhite
	noiseTypeBrown
)

type noiseOsc struct {
	info           core.ModuleInfo // module info
	ntype          noiseType       // noise type
	rand           *core.Rand32    // random state
	b0, b1, b2, b3 float32         // state variables
	b4, b5, b6     float32         // state variables
}

func newNoise(s *core.Synth, ntype noiseType) core.Module {
	m := &noiseOsc{
		info:  noiseOscInfo,
		ntype: ntype,
		rand:  core.NewRand32(0),
	}
	return s.Register(m)
}

// NewNoiseWhite returns a white noise generator module.
// white noise (spectral density = k)
func NewNoiseWhite(s *core.Synth) core.Module {
	log.Info.Printf("")
	return newNoise(s, noiseTypeWhite)
}

// NewNoiseBrown returns a brown noise generator module.
// brown noise (spectral density = k/f*f)
func NewNoiseBrown(s *core.Synth) core.Module {
	log.Info.Printf("")
	return newNoise(s, noiseTypeBrown)
}

// NewNoisePink1 returns a pink noise generator module.
// pink noise (spectral density = k/f): fast, inaccurate version
func NewNoisePink1(s *core.Synth) core.Module {
	log.Info.Printf("")
	return newNoise(s, noiseTypePink1)
}

// NewNoisePink2 returns a pink noise generator module.
// pink noise (spectral density = k/f): slow, accurate version
func NewNoisePink2(s *core.Synth) core.Module {
	log.Info.Printf("")
	return newNoise(s, noiseTypePink2)
}

// Return the child modules.
func (m *noiseOsc) Child() []core.Module {
	return nil
}

// Stop and performs any cleanup of a module.
func (m *noiseOsc) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

//-----------------------------------------------------------------------------

func (m *noiseOsc) generateWhite(out *core.Buf) {
	for i := 0; i < len(out); i++ {
		out[i] = m.rand.Float32()
	}
}

func (m *noiseOsc) generateBrown(out *core.Buf) {
	b0 := m.b0
	for i := 0; i < len(out); i++ {
		white := m.rand.Float32()
		b0 = (b0 + (0.02 * white)) * (1.0 / 1.02)
		out[i] = b0 * (1.0 / 0.38)
	}
	m.b0 = b0
}

func (m *noiseOsc) generatePink1(out *core.Buf) {
	b0 := m.b0
	b1 := m.b1
	b2 := m.b2
	for i := 0; i < len(out); i++ {
		white := m.rand.Float32()
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

func (m *noiseOsc) generatePink2(out *core.Buf) {
	b0 := m.b0
	b1 := m.b1
	b2 := m.b2
	b3 := m.b3
	b4 := m.b4
	b5 := m.b5
	b6 := m.b6
	for i := 0; i < len(out); i++ {
		white := m.rand.Float32()
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
func (m *noiseOsc) Process(buf ...*core.Buf) bool {
	out := buf[0]
	switch m.ntype {
	case noiseTypeWhite:
		m.generateWhite(out)
	case noiseTypePink1:
		m.generatePink1(out)
	case noiseTypePink2:
		m.generatePink2(out)
	case noiseTypeBrown:
		m.generateBrown(out)
	default:
		panic(fmt.Sprintf("bad noise type %d", m.ntype))
	}
	return true
}

//-----------------------------------------------------------------------------
