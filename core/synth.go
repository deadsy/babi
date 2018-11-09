//-----------------------------------------------------------------------------
/*

Synth

*/
//-----------------------------------------------------------------------------

package core

import (
	"errors"

	"github.com/deadsy/babi/utils/cbuf"
	"github.com/deadsy/babi/utils/log"
)

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
	jack  *Jack                // jack client object
	out   []Buf                // audio output buffers
	in    []Buf                // audio input buffers
	event *cbuf.CircularBuffer // event buffer
}

// NewSynth creates a synthesizer object.
func NewSynth() *Synth {
	log.Info.Printf("")
	return &Synth{
		event: cbuf.NewCircularBuffer(numEvents),
	}
}

// SetPatch sets the root module of the synthesizer.
func (s *Synth) SetPatch(m Module) {
	log.Info.Printf(ModuleString(m))
	s.root = m
}

// StartJack starts the jack client.
func (s *Synth) StartJack() error {

	if s.root == nil {
		return errors.New("no root module defined")
	}

	mi := s.root.Info()
	var n int

	// allocate audio input buffers
	n = mi.In.numPorts(PortTypeAudioBuffer)
	if n != 0 {
		s.in = make([]Buf, n)
	}
	// allocate audio output buffers
	n = mi.Out.numPorts(PortTypeAudioBuffer)
	if n != 0 {
		s.out = make([]Buf, n)
	}

	// create the jack client
	jack, err := NewJack(s.root)
	if err != nil {
		return err
	}
	s.jack = jack

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
		for i := 0; i < len(s.out); i++ {
			s.out[i].Zero()
		}
		// process the root module
		if s.root != nil && s.root.Active() {
			s.root.Process(&s.out[0], &s.out[1])
		}
		// write the output to the audio device
		s.jack.WriteAudio(s.out)
	}
}

// Close handles synth cleanup.
func (s *Synth) Close() {
	log.Info.Printf("")
	if s.jack != nil {
		s.jack.Close()
	}
}

//-----------------------------------------------------------------------------
