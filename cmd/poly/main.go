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
	"github.com/deadsy/babi/module/osc"
	"github.com/deadsy/babi/module/patch"
	"github.com/deadsy/babi/module/voice"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

func main() {

	s := core.NewSynth()

	// Pick a voice
	//v := func(s *core.Synth) core.Module { return voice.NewOsc(s, osc.NewSine(s)) }
	//v := func(s *core.Synth) core.Module { return voice.NewOsc(s, osc.NewSquareBasic(s)) }
	//v := func(s *core.Synth) core.Module { return voice.NewOsc(s, osc.NewNoisePink2(s)) }
	//v := func(s *core.Synth) core.Module { return voice.NewOsc(s, osc.NewSawtoothBasic(s)) }
	v := func(s *core.Synth) core.Module { return voice.NewOsc(s, osc.NewGoom(s)) }
	//v := voice.NewKarplusStrong

	// create the polyphonic patch
	p := patch.NewPoly(s, v)

	// set the root patch for the synth
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
