//-----------------------------------------------------------------------------
/*

Time Base

Return an output buffer with the time of each sample.

*/
//-----------------------------------------------------------------------------

package plot

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

// Info returns the module information.
func (m *timeModule) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "time",
		In:   nil,
		Out:  nil,
	}
}

//-----------------------------------------------------------------------------

type timeModule struct {
	synth *core.Synth // top-level synth
	x     uint64      // current x-value
}

// NewTime returns a time base module.
func NewTime(s *core.Synth) core.Module {
	log.Info.Printf("")
	return &timeModule{
		synth: s,
	}
}

// Child returns the child modules of this module.
func (m *timeModule) Child() []core.Module {
	return nil
}

// Stop performs any cleanup of a module.
func (m *timeModule) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *timeModule) Process(buf ...*core.Buf) {
	out := buf[0]
	for i := range out {
		out[i] = float32(m.x) * core.AudioSamplePeriod
		m.x++
	}
}

// Active returns true if the module has non-zero output.
func (m *timeModule) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
