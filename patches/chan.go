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

type channelPatch struct {
	channel []core.Patch
}

func NewChannelPatch(channel []core.Patch) core.Patch {
	log.Info.Printf("")
	return &channelPatch{
		channel: channel,
	}
}

//-----------------------------------------------------------------------------

// Run the patch.
// N buffers in, N buffers out where N is then number of channels.
func (p *channelPatch) Process(in, out []*core.Buf) {
	for i, p := range p.channel {
		if p != nil && p.Active() {
			p.Process(in[i:i+1], out[i:i+1])
		}
	}
}

// Process a patch event.
func (p *channelPatch) Event(e *core.Event) {
	switch e.GetType() {
	case core.Event_MIDI:
		// send the event to the channel patch
		ch := e.GetMIDIEvent().GetChannel()
		if int(ch) < len(p.channel) && p.channel[ch] != nil {
			// send the event to the subpatch
			p.channel[ch].Event(e)
		} else {
			log.Info.Printf("no patch on channel %d", ch)
		}
	default:
		log.Info.Printf("unhandled event %s", e)
	}
}

// Return true if the patch has non-zero output.
func (p *channelPatch) Active() bool {
	return true
}

func (p *channelPatch) Stop() {
	log.Info.Printf("")
	// do nothing
}

//-----------------------------------------------------------------------------
