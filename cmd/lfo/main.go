//-----------------------------------------------------------------------------
/*

LFO Test

*/
//-----------------------------------------------------------------------------

package main

import (
	"os"
	"os/signal"

	"github.com/deadsy/babi/cmd/lfo/app"
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

func main() {

	s := core.NewSynth()

	// create the application patch
	p := app.NewPatch(s, 0)

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
