//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package main

import (
	"fmt"

	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/audio"
	"github.com/deadsy/babi/module/midi"
	"github.com/deadsy/babi/module/noise"
	"github.com/deadsy/babi/module/osc"
	"github.com/deadsy/babi/patches"
)

//-----------------------------------------------------------------------------

func main() {

	// setup audio output
	au, err := core.NewPulse()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	defer au.Close()

	// create the synth

	b := midi.NewPoly(patches.NewSimple, 16)
	_ = b

	d := osc.NewKarplusStrong()
	_ = d

	e := noise.NewWhite()
	_ = e

	f := patches.NewSimple()
	_ = f

	g := audio.NewPan()
	_ = g

	s := core.NewSynth(b, au)
	_ = s

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
