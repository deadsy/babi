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
func (p *channelPatch) Process(in, out []*core.Buf) {
	for i, p := range p.channel {
		if p != nil && p.Active() {
			p.Process(nil, []*core.Buf{out[i]})
		}
	}
}

// Process a patch event.
func (p *channelPatch) Event(e *core.Event) {
	switch e.Etype {
	case core.Event_MIDI:
		// send the event to the channel patch
		ch := e.Info.(*core.MIDIEvent).GetChannel()
		if int(ch) < len(p.channel) && p.channel[ch] != nil {
			// send the event to the subpatch
			p.channel[ch].Event(e)
		} else {
			log.Info.Printf("no patch on channel %d for midi event", ch)
		}
	default:
		log.Info.Printf("unhandled event type %s", e)
	}
}

// Return true if the patch has non-zero output.
func (p *channelPatch) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
