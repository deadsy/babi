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

	s := core.NewSynth(audio)
	s.SetPatch(patches.NewKarplusStrongPatch(s))
	//s.SetPatch(patches.NewSimplePatch(s))
	//s.SetPatch(patches.NewNoisePatch(s))
	s.Run()
}

//-----------------------------------------------------------------------------
