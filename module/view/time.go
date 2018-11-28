//-----------------------------------------------------------------------------
/*

Time Base

Return an output buffer whose values are the sample times.

*/
//-----------------------------------------------------------------------------

package view

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

var timeViewInfo = core.ModuleInfo{
	Name: "timeView",
	In:   nil,
	Out:  nil,
}

// Info returns the module information.
func (m *timeView) Info() *core.ModuleInfo {
	return &timeViewInfo
}

// ID returns the unique module identifier.
func (m *timeView) ID() string {
	return m.id
}

//-----------------------------------------------------------------------------

type timeView struct {
	synth *core.Synth // top-level synth
	id    string      // module identifier
	x     uint64      // current x-value
}

// NewTime returns a time base module.
func NewTime(s *core.Synth) core.Module {
	log.Info.Printf("")
	return &timeView{
		synth: s,
		id:    core.GenerateID(timeViewInfo.Name),
	}
}

// Child returns the child modules of this module.
func (m *timeView) Child() []core.Module {
	return nil
}

// Stop performs any cleanup of a module.
func (m *timeView) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *timeView) Process(buf ...*core.Buf) {
	out := buf[0]
	for i := range out {
		out[i] = float32(m.x) * core.AudioSamplePeriod
		m.x++
	}
}

// Active returns true if the module has non-zero output.
func (m *timeView) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
