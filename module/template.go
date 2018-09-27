//-----------------------------------------------------------------------------
/*

Module Name and Description

*/
//-----------------------------------------------------------------------------

package module

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

// Info returns the module information.
func (m *xModule) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "x",
		In:   nil,
		Out:  nil,
	}
}

//-----------------------------------------------------------------------------

type xModule struct {
	synth *core.Synth // top-level synth
}

// NewX returns an X module.
func NewX(s *core.Synth) core.Module {
	log.Info.Printf("")
	return &xModule{
		synth: s,
	}
}

// Child returns the child modules of this module.
func (m *xModule) Child() []core.Module {
	return nil
}

// Stop performs any cleanup of a module.
func (m *xModule) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *xModule) Process(buf ...*core.Buf) {
}

// Active returns true if the module has non-zero output.
func (m *xModule) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
