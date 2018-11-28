//-----------------------------------------------------------------------------
/*

Goom Voice Control Module

A goom voice has more controls than I have knobs on my MIDI controller.
This module alows modal switching between the 8 knobs I do have.
That is: Hit a drum pad, switch modes to a different control group.

*/
//-----------------------------------------------------------------------------

package goom

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

var ctrlGoomInfo = core.ModuleInfo{
	Name: "ctrlGoom",
	In:   nil,
	Out:  nil,
}

// Info returns the module information.
func (m *ctrlGoom) Info() *core.ModuleInfo {
	return &ctrlGoomInfo
}

// ID returns the unique module identifier.
func (m *ctrlGoom) ID() string {
	return m.id
}

//-----------------------------------------------------------------------------

type ctrlGoom struct {
	synth *core.Synth // top-level synth
	id    string      // module identifier
}

// NewX returns an X module.
func NewX(s *core.Synth) core.Module {
	log.Info.Printf("")
	return &ctrlGoom{
		synth: s,
		id:    core.GenerateID(ctrlGoomInfo.Name),
	}
}

// Child returns the child modules of this module.
func (m *ctrlGoom) Child() []core.Module {
	return nil
}

// Stop performs any cleanup of a module.
func (m *ctrlGoom) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *ctrlGoom) Process(buf ...*core.Buf) {
}

// Active returns true if the module has non-zero output.
func (m *ctrlGoom) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
