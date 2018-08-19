//-----------------------------------------------------------------------------
/*

Channel Patch

Demux incoming events into N subpatches based on channel number.

*/
//-----------------------------------------------------------------------------

package patches

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/log"
)

//-----------------------------------------------------------------------------

const MAX_CHANNELS = 16

//-----------------------------------------------------------------------------

type channelPatch struct {
	channel *[]core.Patch
}

func NewChannelPatch(channel *[]core.Patch) core.Patch {
	log.Info.Printf("")
	return &channelPatch{
		channel: channel,
	}
}

//-----------------------------------------------------------------------------

// Run the patch.
func (p *channelPatch) Process() {
	for _, p := range *p.channel {
		if p != nil && p.Active() {
			p.Process()
		}
	}
}

// Process a patch event.
func (p *channelPatch) Event(e *core.Event) {
}

// Return true if the patch has non-zero output.
func (p *channelPatch) Active() bool {
	return true
}

// Output to the parent patch.
func (p *channelPatch) Out(out ...*core.Buf) {
}

//-----------------------------------------------------------------------------
