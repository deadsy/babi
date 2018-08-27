//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package main

import (
	"fmt"

	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/env"
	"github.com/deadsy/babi/module/midi"
	"github.com/deadsy/babi/module/osc"
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
	_ = s

	a := env.NewADSR()
	_ = a

	b := midi.NewPoly(env.NewADSR, 16)
	_ = b

	c := osc.NewSine()
	_ = c

	d := osc.NewKarplusStrong()
	_ = d

	// 	// create the patches
	// 	//p0 := patches.NewPolyPatch(patches.NewSimplePatch, 16)
	// 	p0 := patches.NewPolyPatch(patches.NewKSPatch, 16)
	// 	p1 := patches.NewChannelPatch([]core.Patch{p0})
	//
	// 	// set the root patch and run the synth
	// 	s.SetRoot(p1)
	// 	s.Run()
}

//-----------------------------------------------------------------------------
