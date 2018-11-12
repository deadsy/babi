//-----------------------------------------------------------------------------
/*

Jack Client Object

*/
//-----------------------------------------------------------------------------

package core

import (
	"errors"
	"fmt"

	"github.com/deadsy/babi/jack"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------
// jack server callbacks

func (j *Jack) process(nframes uint32) int {
	//log.Info.Printf("")

	// process MIDI input events
	for i := range j.midiIn {
		p := &j.midiIn[i]
		event := p.port.GetMIDIEvents(nframes)
		for j := range event {
			e := &event[j]
			midiEvent := convertToMIDIEvent(e.Data)
			if midiEvent != nil {
				log.Info.Printf("%s", midiEvent.String())
				//s.PushEvent(nil, "midi_in", midiEvent)
			}
		}
	}

	// get audio output buffers
	for i := range j.audioOut {
		j.audioOut[i].audio = j.audioOut[i].port.GetBuffer(nframes)
	}

	// get audio input buffers
	for i := range j.audioIn {
		j.audioIn[i].audio = j.audioIn[i].port.GetBuffer(nframes)
	}

	return 0
}

func (j *Jack) shutdown() {
	log.Info.Printf("")
}

//-----------------------------------------------------------------------------

type jackPort struct {
	port  *jack.Port         // jack port
	audio []jack.AudioSample // audio sample buffer
}

// Jack contains state for the jack client.
type Jack struct {
	name     string // client name
	client   *jack.Client
	audioOut []jackPort // audio output ports
	audioIn  []jackPort // audio input ports
	midiOut  []jackPort // midi output ports
	midiIn   []jackPort // midi input ports
}

// NewJack returns a jack client object.
func NewJack(name string, m Module) (*Jack, error) {
	j := &Jack{name: name}

	log.Info.Printf("jack version %s", jack.GetVersionString())

	// open the client
	client, status := jack.ClientOpen(name, jack.NoStartServer)
	if status != 0 {
		j.Close()
		return nil, errors.New(status.String())
	}
	j.client = client

	// check sample rate
	rate := client.GetSampleRate()
	if rate != AudioSampleFrequency {
		j.Close()
		return nil, fmt.Errorf("jack sample rate %d != babi sample rate %d", rate, AudioSampleFrequency)
	}

	// check audio buffer size
	bufsize := client.GetBufferSize()
	if bufsize != AudioBufferSize {
		j.Close()
		return nil, fmt.Errorf("jack buffer size %d != babi buffer size %d", bufsize, AudioBufferSize)
	}

	// tell the JACK server to call process() whenever there is work to be done.
	rc := client.SetProcessCallback(func(nframes uint32) int { return j.process(nframes) })
	if rc != 0 {
		j.Close()
		return nil, fmt.Errorf("SetProcessCallback() error %d", rc)
	}

	// tell the JACK server to call shutdown() if it ever shuts down,
	// either entirely, or if it just decides to stop calling us.
	client.OnShutdown(func() { j.shutdown() })

	mi := m.Info()

	// audio output ports
	n := mi.Out.numPorts(PortTypeAudioBuffer)
	ports, err := j.registerPorts(n, "audio_out", jack.DefaultAudio, jack.PortIsOutput)
	if err != nil {
		j.Close()
		return nil, err
	}
	j.audioOut = ports

	// audio input ports
	n = mi.In.numPorts(PortTypeAudioBuffer)
	ports, err = j.registerPorts(n, "audio_in", jack.DefaultAudio, jack.PortIsInput)
	if err != nil {
		j.Close()
		return nil, err
	}
	j.audioIn = ports

	// MIDI output ports
	n = mi.Out.numPorts(PortTypeMIDI)
	ports, err = j.registerPorts(n, "midi_out", jack.DefaultMIDI, jack.PortIsOutput)
	if err != nil {
		j.Close()
		return nil, err
	}
	j.midiOut = ports

	// MIDI input ports
	n = mi.In.numPorts(PortTypeMIDI)
	ports, err = j.registerPorts(n, "midi_in", jack.DefaultMIDI, jack.PortIsInput)
	if err != nil {
		j.Close()
		return nil, err
	}
	j.midiIn = ports

	// Tell the JACK server that we are ready to roll.
	// Our process() callback will start running now.
	rc = client.Activate()
	if rc != 0 {
		j.Close()
		return nil, fmt.Errorf("Activate() error %d", rc)
	}

	return j, nil
}

// Close closes the jack client.
func (j *Jack) Close() {
	log.Info.Printf("")
	if j.client == nil {
		return
	}
	j.client.Deactivate()
	j.unregisterPorts(j.audioOut)
	j.unregisterPorts(j.audioIn)
	j.unregisterPorts(j.midiOut)
	j.unregisterPorts(j.midiIn)
	j.client.Close()
	j.client = nil
}

// WriteAudio writes data to an audio stream.
func (j *Jack) WriteAudio(audio []Buf) {
}

//-----------------------------------------------------------------------------

// registerPorts registers a set of jack client ports.
func (j *Jack) registerPorts(n int, prefix, portType string, flags uint64) ([]jackPort, error) {
	ports := make([]jackPort, n)
	for i := range ports {
		pname := fmt.Sprintf("%s_%d", prefix, i)
		p := j.client.PortRegister(pname, portType, flags, 0)
		if p == nil {
			return nil, fmt.Errorf("can't register port %s:%s", j.name, pname)
		}
		log.Info.Printf("registered \"%s\"", p.Name())
		ports[i].port = p
	}
	return ports, nil
}

// unregisterPorts unregisters a set of jack clientports.
func (j *Jack) unregisterPorts(ports []jackPort) {
	if j.client == nil || ports == nil {
		return
	}
	for _, p := range ports {
		log.Info.Printf("unregistered \"%s\"", p.port.Name())
		j.client.PortUnregister(p.port)
	}
}

//-----------------------------------------------------------------------------
