//-----------------------------------------------------------------------------
/*

Noise Objects

https://noisehack.com/generate-noise-web-audio-api/
http://www.musicdsp.org/files/pink.txt
https://en.wikipedia.org/wiki/Pink_noise
https://en.wikipedia.org/wiki/White_noise
https://en.wikipedia.org/wiki/Brownian_noise

*/
//-----------------------------------------------------------------------------

package noise

import "github.com/deadsy/babi/core"

//-----------------------------------------------------------------------------

type White struct {
	r *core.Rand
}

func NewWhite() *White {
	return &White{
		r: core.NewRand(0),
	}
}

// white noise (spectral density = k)
func (n *White) Process(out *core.SBuf) {
	for i := 0; i < len(out); i++ {
		out[i] = n.r.Float()
	}
}

//-----------------------------------------------------------------------------

type Brown struct {
	r  *core.Rand
	b0 float32
}

func NewBrown() *Brown {
	return &Brown{
		r: core.NewRand(0),
	}
}

// brown noise (spectral density = k/f*f)
func (n *Brown) Process(out *core.SBuf) {
	b0 := n.b0
	for i := 0; i < len(out); i++ {
		white := n.r.Float()
		b0 = (b0 + (0.02 * white)) * (1.0 / 1.02)
		out[i] = b0 * (1.0 / 0.38)
	}
	n.b0 = b0
}

//-----------------------------------------------------------------------------

type Pink1 struct {
	r          *core.Rand
	b0, b1, b2 float32
}

func NewPink1() *Pink1 {
	return &Pink1{
		r: core.NewRand(0),
	}
}

// pink noise (spectral density = k/f): fast, inaccurate version
func (n *Pink1) Process(out *core.SBuf) {
	b0 := n.b0
	b1 := n.b1
	b2 := n.b2
	for i := 0; i < len(out); i++ {
		white := n.r.Float()
		b0 = 0.99765*b0 + white*0.0990460
		b1 = 0.96300*b1 + white*0.2965164
		b2 = 0.57000*b2 + white*1.0526913
		pink := b0 + b1 + b2 + white*0.1848
		out[i] = pink * (1.0 / 10.4)
	}
	n.b0 = b0
	n.b1 = b1
	n.b2 = b2
}

//-----------------------------------------------------------------------------

type Pink2 struct {
	r                          *core.Rand
	b0, b1, b2, b3, b4, b5, b6 float32
}

func NewPink2() *Pink2 {
	return &Pink2{
		r: core.NewRand(0),
	}
}

// pink noise (spectral density = k/f): slow, accurate version
func (n *Pink2) Process(out *core.SBuf) {
	b0 := n.b0
	b1 := n.b1
	b2 := n.b2
	b3 := n.b3
	b4 := n.b4
	b5 := n.b5
	b6 := n.b6
	for i := 0; i < len(out); i++ {
		white := n.r.Float()
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
	n.b0 = b0
	n.b1 = b1
	n.b2 = b2
	n.b3 = b3
	n.b4 = b4
	n.b5 = b5
	n.b6 = b6
}

//-----------------------------------------------------------------------------
