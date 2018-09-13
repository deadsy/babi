//-----------------------------------------------------------------------------
/*

Audio Interface Objects

*/
//-----------------------------------------------------------------------------

package core

import (
	"errors"

	"github.com/deadsy/babi/pulse"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

type Audio interface {
	Close()
	Write(l, r *Buf)
}

//-----------------------------------------------------------------------------
// pulse audio

type Pulse struct {
	pa  *pulse.PulseMainLoop
	ctx *pulse.PulseContext
	st  *pulse.PulseStream
}

func NewPulse() (Audio, error) {
	log.Info.Printf("")

	pa := pulse.NewPulseMainLoop()
	pa.Start()

	ctx := pa.NewContext("default", 0)
	if ctx == nil {
		pa.Dispose()
		return nil, errors.New("failed to create a new context")
	}

	st := ctx.NewStream("default", &pulse.PulseSampleSpec{Format: pulse.SAMPLE_FLOAT32LE, Rate: AudioSampleFrequency, Channels: 2})
	if st == nil {
		ctx.Dispose()
		pa.Dispose()
		return nil, errors.New("failed to create a new stream")
	}
	st.ConnectToSink()

	return &Pulse{pa, ctx, st}, nil
}

func (p *Pulse) Close() {
	log.Info.Printf("")
	p.st.Dispose()
	p.ctx.Dispose()
	p.pa.Dispose()
}

func (p *Pulse) Write(l, r *Buf) {
	// combine left/right channels into a single slice.
	buf := make([]float32, 2*AudioBufferSize)
	for i := 0; i < AudioBufferSize; i += 1 {
		buf[2*i] = l[i]
		buf[(2*i)+1] = r[i]
	}
	p.st.Write(buf, pulse.SEEK_RELATIVE)
}

//-----------------------------------------------------------------------------
