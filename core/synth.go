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

type Synth struct {
	root  Module              // root module
	audio Audio               // audio output device
	out   [AUDIO_CHANNELS]Buf // audio output buffers
	in    [AUDIO_CHANNELS]Buf // audio input buffers
}

// NewSynth creates a synthesizer object.
func NewSynth(root Module, audio Audio) *Synth {
	log.Info.Printf("")
	return &Synth{
		root:  root,
		audio: audio,
	}
}

// Main loop for the synthesizer.
func (s *Synth) Run() {

	s.root.Event(NewEventMIDI(EventMIDI_NoteOn, 0, 69, 127))

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
			s.root.Process(&s.in[0], &s.in[1], &s.out[0], &s.out[1])
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
