//-----------------------------------------------------------------------------
/*

Audio

*/
//-----------------------------------------------------------------------------

package core

import (
	"errors"

	"github.com/deadsy/babi/pulse"
)

//-----------------------------------------------------------------------------

type Audio interface {
	Close()
	Write(l, r *Buf)
}

//-----------------------------------------------------------------------------

type Pulse struct {
	pa  *pulse.PulseMainLoop
	ctx *pulse.PulseContext
	st  *pulse.PulseStream
}

func NewPulse() (Audio, error) {

	pa := pulse.NewPulseMainLoop()
	pa.Start()

	ctx := pa.NewContext("default", 0)
	if ctx == nil {
		pa.Dispose()
		return nil, errors.New("failed to create a new context")
	}

	st := ctx.NewStream("default", &pulse.PulseSampleSpec{Format: pulse.SAMPLE_FLOAT32LE, Rate: AUDIO_FS, Channels: 1})
	if st == nil {
		ctx.Dispose()
		pa.Dispose()
		return nil, errors.New("failed to create a new stream")
	}
	st.ConnectToSink()

	return &Pulse{pa, ctx, st}, nil
}

func (p *Pulse) Close() {
	p.st.Dispose()
	p.ctx.Dispose()
	p.pa.Dispose()
}

func (p *Pulse) Write(l, r *Buf) {
	p.st.Write(l[:], pulse.SEEK_RELATIVE)
}

//-----------------------------------------------------------------------------

type Babi struct {
	patch        Patch
	audio        Audio
	out_l, out_r Buf
}

func NewBabi(audio Audio) *Babi {
	return &Babi{
		audio: audio,
	}
}

func (b *Babi) AddPatch(patch Patch) {
	b.patch = patch
}

func (b *Babi) Run() {
	for {
		b.OutZero()
		b.patch.Process()
		b.audio.Write(&b.out_l, &b.out_r)
	}
}

func (b *Babi) OutZero() {
	b.out_l.Zero()
	b.out_r.Zero()
}

func (b *Babi) OutLR(l, r *Buf) {
	b.out_l.Add(l)
	b.out_r.Add(r)
}

//-----------------------------------------------------------------------------
