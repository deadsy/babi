//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package main

import (
	"os"
	"os/signal"

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

// signalHandler waits for ctrl-C.
func signalHandler(signals chan os.Signal, done chan bool) {
	<-signals
	done <- true
}

//-----------------------------------------------------------------------------

func main() {

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
		s.Close()
		os.Exit(1)
	}

	// setup signal handling
	done := make(chan bool)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go signalHandler(signals, done)
	<-done

	s.Close()
	os.Exit(0)
}

//-----------------------------------------------------------------------------
