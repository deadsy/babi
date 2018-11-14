//-----------------------------------------------------------------------------
/*

Polyphonic Voice Player

*/
//-----------------------------------------------------------------------------

package main

import (
	"os"
	"os/signal"

	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/patch"
	"github.com/deadsy/babi/module/voice"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

func main() {

	s := core.NewSynth()

	//p := patch.NewBasicPatch(s, osc.NewSine(s))
	//p := patch.NewBasicPatch(s, osc.NewSquareBasic(s))
	//p := patch.NewBasicPatch(s, osc.NewNoisePink2(s))
	//p := patch.NewBasicPatch(s, osc.NewSawtoothBasic(s))
	//p := patch.NewBasicPatch(s, osc.NewGoom(s))
	//p := patch.NewKarplusStrongPatch(s)
	//p := patch.NewSequencerTest(s, metronome)

	v := voice.NewKarplusStrong
	p := patch.NewPoly(s, v)

	// set the root patch
	s.SetPatch(p)

	// start the jack client
	err := s.StartJack("babi")
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
