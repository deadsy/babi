//-----------------------------------------------------------------------------
/*

Synth

*/
//-----------------------------------------------------------------------------

package core

import (
	"errors"
	"fmt"

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

// pushEvent pushes an event onto the synth event queue.
func (s *Synth) pushEvent(m Module, name string, e *Event) {
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
	s.nIn = mi.In.numPortsByType(PortTypeAudio)
	s.nOut = mi.Out.numPortsByType(PortTypeAudio)
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
		EventIn(e.dst, e.port, e.event)
	}
	// zero the audio output buffers
	for i := s.nIn; i < len(s.audio); i++ {
		s.audio[i].Zero()
	}
	// process the root module
	if s.root != nil {
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

// Register registers a new module with the synth.
func (s *Synth) Register(m Module) Module {
	mi := m.Info()
	// set a reference to the top-level synth in the module info
	mi.Synth = s
	// build the name to port function mapping for the inputs
	mi.inMap = make(map[string]PortFuncType)
	for i := range mi.In {
		name := mi.In[i].Name
		if _, ok := mi.inMap[name]; ok {
			panic(fmt.Sprintf("module \"%s\" must have only one input port with name \"%s\"", mi.Name, name))
		}
		mi.inMap[name] = mi.In[i].PortFunc
	}
	// build the name to port mapping for the outputs
	mi.outMap = make(map[string][]dstPort)
	for i := range mi.Out {
		name := mi.Out[i].Name
		if _, ok := mi.outMap[name]; ok {
			panic(fmt.Sprintf("module \"%s\" must have only one output port with name \"%s\"", mi.Name, name))
		}
		mi.outMap[name] = nil
	}
	return m
}

//-----------------------------------------------------------------------------
