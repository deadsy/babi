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

	m0 := patches.NewSimple()

	// create the synth
	s := core.NewSynth(m0, audio)
	s.Run()
}

//-----------------------------------------------------------------------------
