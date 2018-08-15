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

	// setup synth and add patches
	s := core.NewSynth(audio)
	s.AddPatch(&patches.KarplusStrongInfo, 0)
	s.AddPatch(&patches.SimpleInfo, 1)
	s.VoiceAlloc(0, 69)

	// run the synth
	s.Run()
}

//-----------------------------------------------------------------------------
