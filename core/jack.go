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

func processCallback(nframes uint32) int {
	log.Info.Printf("")
	return 0
}

func shutdown() {
	log.Info.Printf("")
}

func infoShutdown() {
	log.Info.Printf("")
}

//-----------------------------------------------------------------------------

// Jack contains state for the jack client.
type Jack struct {
	name         string // client name
	client       *jack.Client
	audioPortOut []*jack.Port // audio output ports
	audioPortIn  []*jack.Port // audio input ports
	midiPortOut  []*jack.Port // midi output ports
	midiPortIn   []*jack.Port // midi input ports
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

	mi := m.Info()

	// audio output ports
	n := mi.Out.numPorts(PortTypeAudioBuffer)
	ports, err := j.registerPorts(n, "audio_out", jack.DefaultAudio, jack.PortIsOutput)
	if err != nil {
		j.Close()
		return nil, err
	}
	j.audioPortOut = ports

	// audio input ports
	n = mi.In.numPorts(PortTypeAudioBuffer)
	ports, err = j.registerPorts(n, "audio_in", jack.DefaultAudio, jack.PortIsInput)
	if err != nil {
		j.Close()
		return nil, err
	}
	j.audioPortIn = ports

	// midi output ports
	n = mi.Out.numPorts(PortTypeMIDI)
	ports, err = j.registerPorts(n, "midi_out", jack.DefaultMIDI, jack.PortIsOutput)
	if err != nil {
		j.Close()
		return nil, err
	}
	j.midiPortOut = ports

	// midi input ports
	n = mi.In.numPorts(PortTypeMIDI)
	ports, err = j.registerPorts(n, "midi_in", jack.DefaultMIDI, jack.PortIsInput)
	if err != nil {
		j.Close()
		return nil, err
	}
	j.midiPortIn = ports

	rc := client.SetProcessCallback(processCallback)
	if rc != 0 {
		j.Close()
		return nil, fmt.Errorf("SetProcessCallback() error %d", rc)
	}

	client.OnShutdown(shutdown)
	client.OnInfoShutdown(infoShutdown)

	return j, nil
}

// Close closes the jack client.
func (j *Jack) Close() {
	log.Info.Printf("")

	j.unregisterPorts(j.audioPortOut)
	j.audioPortOut = nil

	j.unregisterPorts(j.audioPortIn)
	j.audioPortIn = nil

	j.unregisterPorts(j.midiPortOut)
	j.midiPortOut = nil

	j.unregisterPorts(j.midiPortIn)
	j.midiPortIn = nil

	if j.client != nil {
		j.client.Close()
		j.client = nil
	}
}

// WriteAudio writes data to an audio stream.
func (j *Jack) WriteAudio(audio []Buf) {
}

//-----------------------------------------------------------------------------

// registerPorts registers a set of jack client ports.
func (j *Jack) registerPorts(n int, prefix, portType string, flags uint64) ([]*jack.Port, error) {
	ports := make([]*jack.Port, n)
	for i := range ports {
		pname := fmt.Sprintf("%s_%d", prefix, i+1)
		p := j.client.PortRegister(pname, portType, flags, 0)
		if p == nil {
			return nil, fmt.Errorf("can't register port %s:%s", j.name, pname)
		}
		log.Info.Printf("registered \"%s\"", p.Name())
		ports[i] = p
	}
	return ports, nil
}

// unregisterPorts unregisters a set of jack client ports.
func (j *Jack) unregisterPorts(ports []*jack.Port) {
	if j.client == nil || ports == nil {
		return
	}
	for i, p := range ports {
		log.Info.Printf("unregistered \"%s\"", p.Name())
		j.client.PortUnregister(p)
		ports[i] = nil
	}
}

//-----------------------------------------------------------------------------
