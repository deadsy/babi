//-----------------------------------------------------------------------------
/*

Patch Operations

*/
//-----------------------------------------------------------------------------

package core

import (
	"errors"
)

//-----------------------------------------------------------------------------

type Patch interface {
	Process()
	Active() bool
}

type PatchInfo struct {
	Name string
	New  func() Patch
}

//-----------------------------------------------------------------------------

var channel_to_patch [16]*PatchInfo

// RegisterPatch assigns a patch to a midi channel
func RegisterPatch(p *PatchInfo, ch int) error {
	if ch < 0 || ch >= len(channel_to_patch) {
		return errors.New("channel value is out of range")
	}
	channel_to_patch[ch] = p
	return nil
}

//-----------------------------------------------------------------------------
