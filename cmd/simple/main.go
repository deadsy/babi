//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package main

import (
	"os"

	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/osc"
	"github.com/deadsy/babi/module/seq"
	"github.com/deadsy/babi/patches"
	"github.com/deadsy/babi/utils/log"
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
	rc := 0

	s := core.NewSynth()

	//p := patches.NewBasicPatch(s, osc.NewSine(s))
	//p := patches.NewBasicPatch(s, osc.NewSquareBasic(s))
	//p := patches.NewBasicPatch(s, osc.NewNoisePink2(s))
	//p := patches.NewBasicPatch(s, osc.NewSawtoothBasic(s))
	p := patches.NewBasicPatch(s, osc.NewGoom(s))
	//p := patches.NewKarplusStrongPatch(s)
	//p := patches.NewSequencerTest(s, metronome)

	// set the root patch
	s.SetPatch(p)

	// start the jack client
	err := s.StartJack("simple")
	if err != nil {
		log.Error.Printf("%s", err)
		rc = 1
		goto Exit
	}

	// run the synth
	s.Run()

Exit:
	s.Close()
	os.Exit(rc)
}

//-----------------------------------------------------------------------------
