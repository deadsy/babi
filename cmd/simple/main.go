//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package main

import (
	"fmt"

	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/seq"
	"github.com/deadsy/babi/patches"
)

//-----------------------------------------------------------------------------

var metronome = []seq.Op{
	seq.OpNote(1, 69, 100, 4),
	seq.OpRest(12),
	seq.OpNote(1, 60, 100, 4),
	seq.OpRest(12),
	seq.OpNote(1, 60, 100, 4),
	seq.OpRest(12),
	seq.OpNote(1, 60, 100, 4),
	seq.OpRest(12),
	seq.OpLoop(),
}

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
	//s.SetPatch(patches.NewKarplusStrongPatch(s))
	s.SetPatch(patches.NewSequencerTest(s, metronome))
	//s.SetPatch(patches.NewSimplePatch(s))
	//s.SetPatch(patches.NewNoisePatch(s))
	s.Run()
}

//-----------------------------------------------------------------------------
