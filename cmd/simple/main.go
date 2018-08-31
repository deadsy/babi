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
	m1 := patches.NewPan(0, m0)

	// create the synth
	s := core.NewSynth(m1, audio)
	s.Run()
}

//-----------------------------------------------------------------------------
