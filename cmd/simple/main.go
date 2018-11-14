//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package main

import (
	"os"
	"os/signal"

	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/osc"
	"github.com/deadsy/babi/module/patch"
	"github.com/deadsy/babi/module/seq"
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

	s := core.NewSynth()

	//p := patch.NewBasicPatch(s, osc.NewSine(s))
	//p := patch.NewBasicPatch(s, osc.NewSquareBasic(s))
	//p := patch.NewBasicPatch(s, osc.NewNoisePink2(s))
	//p := patch.NewBasicPatch(s, osc.NewSawtoothBasic(s))
	p := patch.NewBasicPatch(s, osc.NewGoom(s))
	//p := patch.NewKarplusStrongPatch(s)
	//p := patch.NewSequencerTest(s, metronome)

	// set the root patch
	s.SetPatch(p)

	// start the jack client
	err := s.StartJack("simple")
	if err != nil {
		log.Error.Printf("%s", err)
		s.Close()
		os.Exit(1)
	}

	// signal handling
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	<-signals

	s.Close()
	os.Exit(0)
}

//-----------------------------------------------------------------------------
