//-----------------------------------------------------------------------------
/*

Left/Right Pan and Volume Patch

*/
//-----------------------------------------------------------------------------

package patches

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/log"
	"github.com/deadsy/babi/objects/audio"
)

//-----------------------------------------------------------------------------

type panPatch struct {
	pan *audio.Pan // left right panning
}

func NewPanPatch() core.Patch {
	log.Info.Printf("")
	p := &panPatch{
		pan: audio.NewPan(),
	}
	p.pan.SetPan(0.5)
	p.pan.SetVol(1.0)
	return p
}

//-----------------------------------------------------------------------------

// Run the patch. len(in) = 1, len(out) = 2
func (p *panPatch) Process(in, out []*core.Buf) {
	// pan input to left/right output channels
	p.pan.Process(in[0], out[0], out[1])
}

// Process a patch event.
func (p *panPatch) Event(e *core.Event) {
}

// Return true if the patch has non-zero output.
func (p *panPatch) Active() bool {
	return true
}

func (p *panPatch) Stop() {
	// do nothing
}

//-----------------------------------------------------------------------------
