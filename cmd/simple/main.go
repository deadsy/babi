//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package main

import (
	"fmt"

	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/patches"
)

//-----------------------------------------------------------------------------

func main() {
	// setup audio output
	audio, err := core.NewPulse()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	defer audio.Close()

	// create the synth
	s := core.NewSynth(audio)

	// create the patches
	p0 := patches.NewPolyPatch(patches.NewSimplePatch, 16)
	p1 := patches.NewChannelPatch([]core.Patch{p0})

	// set the root patch and run the synth
	s.SetRoot(p1)
	s.Run()
}

//-----------------------------------------------------------------------------
