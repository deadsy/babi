//-----------------------------------------------------------------------------
/*

Synth

*/
//-----------------------------------------------------------------------------

package core

import (
	"github.com/deadsy/babi/utils/cbuf"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

const numMIDIIn = 1
const numAudioIn = 0
const numAudioOut = 2
const numEvents = 32

//-----------------------------------------------------------------------------

// Synthesizer
type Synth struct {
	root  Module               // root module
	audio Audio                // audio output device
	out   [numAudioOut]Buf     // audio output buffers
	event *cbuf.CircularBuffer // event buffer
}

// NewSynth creates a synthesizer object.
func NewSynth(audio Audio) *Synth {
	log.Info.Printf("")
	return &Synth{
		audio: audio,
		event: cbuf.NewCircularBuffer(numEvents),
	}
}

// SetPatch sets the root module of the synthesizer.
func (s *Synth) SetPatch(m Module) error {
	err := m.Info().CheckIO(numMIDIIn, numAudioIn, numAudioOut)
	if err != nil {
		return err
	}
	log.Info.Printf(ModuleString(m))
	s.root = m
	return nil
}

// PushEvent pushes an event onto the event queue for future processing.
func (s *Synth) PushEvent(e *Event) {
	err := s.event.Write(e)
	if err != nil {
		log.Info.Printf("%s", err)
	}
}

// Run runs the main loop for the synthesizer.
func (s *Synth) Run() {

	s.PushEvent(NewEventMIDI(EventMIDI_NoteOn, 0, 69, 127))

	for {
		// process all queued events
		for !s.event.Empty() {
			e, _ := s.event.Read()
			s.root.Event(e.(*Event))
		}
		// zero the audio output buffers
		for i := 0; i < numAudioOut; i++ {
			s.out[i].Zero()
		}
		// process the root module
		if s.root != nil && s.root.Active() {
			s.root.Process(&s.out[0], &s.out[1])
		}
		// write the output to the audio device
		s.audio.Write(&s.out[0], &s.out[1])
	}
}

//-----------------------------------------------------------------------------
