//-----------------------------------------------------------------------------
/*

Synth

*/
//-----------------------------------------------------------------------------

package core

import (
	"github.com/deadsy/babi/log"
)

//-----------------------------------------------------------------------------

const AUDIO_CHANNELS = 2

//-----------------------------------------------------------------------------
// patches

type Patch interface {
	Process(in, out []*Buf) // run the patch
	Event(e *Event)         // process an event
	Active() bool           // return true if the patch has non-zero output
	Stop()                  // stop the patch
}

//-----------------------------------------------------------------------------

type Synth struct {
	root  Patch               // root patch
	audio Audio               // audio output device
	out   [AUDIO_CHANNELS]Buf // audio output buffers
	in    [AUDIO_CHANNELS]Buf // audio input buffers
}

// NewSynth creates a synthesizer object.
func NewSynth(audio Audio) *Synth {
	log.Info.Printf("")
	return &Synth{
		audio: audio,
	}
}

// Set the root patch for the synthesizer.
func (s *Synth) SetRoot(p Patch) {
	s.root = p
}

// Main loop for the synthesizer.
func (s *Synth) Run() {

	s.root.Event(NewMIDIEvent(MIDIEvent_NoteOn, 0, 69, 127))

	for {
		// zero the audio output buffers
		for i := 0; i < AUDIO_CHANNELS; i++ {
			s.out[i].Zero()
		}
		// TODO get the audio input buffers
		for i := 0; i < AUDIO_CHANNELS; i++ {
			s.in[i].Zero()
		}
		// process the patches
		if s.root != nil && s.root.Active() {
			// TODO fix buffer handling
			s.root.Process([]*Buf{&s.in[0]}, []*Buf{&s.out[0]})
		}
		// write the output to the audio device
		s.audio.Write(&s.out[0], &s.out[1])
	}
}

func (s *Synth) OutLR(l, r *Buf) {
	s.out[0].Add(l)
	s.out[1].Add(r)
}

//-----------------------------------------------------------------------------
