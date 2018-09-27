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

//-----------------------------------------------------------------------------

const numEvents = 32

// QueueEvent contains an event for future processing.
type QueueEvent struct {
	dst   Module // destination module
	port  string // port name
	event *Event // event
}

// PushEvent pushes an event onto the synth event queue.
func (s *Synth) PushEvent(m Module, name string, e *Event) {
	err := s.event.Write(&QueueEvent{m, name, e})
	if err != nil {
		log.Info.Printf("%s", err)
	}
}

//-----------------------------------------------------------------------------

// Synth is the top-level synthesizer object.
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

// Run runs the main loop for the synthesizer.
func (s *Synth) Run() {
	s.PushEvent(nil, "midi_in", NewEventMIDI(EventMIDINoteOn, 0, 69, 127))
	for {
		// process all queued events
		for !s.event.Empty() {
			x, _ := s.event.Read()
			e := x.(*QueueEvent)
			if e.dst == nil {
				e.dst = s.root
			}
			SendEvent(e.dst, e.port, e.event)
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
