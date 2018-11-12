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
	audio []*Buf               // audio buffers (in + out)
	nIn   int                  // number of audio input buffers
	nOut  int                  // number of audio output buffers
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
	// allocate audio buffers
	mi := m.Info()
	s.nIn = mi.In.numPorts(PortTypeAudioBuffer)
	s.nOut = mi.Out.numPorts(PortTypeAudioBuffer)
	n := s.nIn + s.nOut
	if n != 0 {
		s.audio = make([]*Buf, n)
	}
	for i := range s.audio {
		s.audio[i] = &Buf{}
	}
}

// StartJack starts the jack client.
func (s *Synth) StartJack(name string) error {
	if s.root == nil {
		return errors.New("no root module defined")
	}
	// create the jack client
	jack, err := NewJack(name, s)
	if err != nil {
		return err
	}
	s.jack = jack
	return nil
}

// Loop runs a single iteration of the synthesizer.
func (s *Synth) Loop() {
	// process all queued events
	for !s.event.Empty() {
		x, _ := s.event.Read()

		if x == nil {
			panic("here")
		}

		e := x.(*QueueEvent)
		if e.dst == nil {
			e.dst = s.root
		}
		SendEvent(e.dst, e.port, e.event)
	}
	// zero the audio output buffers
	for i := s.nIn; i < len(s.audio); i++ {
		s.audio[i].Zero()
	}
	// process the root module
	if s.root != nil && s.root.Active() {
		s.root.Process(s.audio...)
	}
}

// Close handles synth cleanup.
func (s *Synth) Close() {
	log.Info.Printf("")
	if s.jack != nil {
		s.jack.Close()
	}
	ModuleStop(s.root)
}

//-----------------------------------------------------------------------------
