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

var xModuleInfo = core.ModuleInfo{
	Name: "xModule",
	In:   nil,
	Out:  nil,
}

// Info returns the general module information.
func (m *xModule) Info() *core.ModuleInfo {
	return &m.info
}

//-----------------------------------------------------------------------------

type xModule struct {
	info core.ModuleInfo // module info
}

// NewX returns an X module.
func NewX(s *core.Synth) core.Module {
	log.Info.Printf("")
	m := &xModule{
		info: xModuleInfo,
	}
	return s.Register(m)
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
