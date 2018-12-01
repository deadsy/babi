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
	//log.Info.Printf("nframes %d", nframes)

	// read MIDI input events
	for _, p := range j.midiIn {
		event := p.GetMIDIEvents(nframes)
		for k := range event {
			e := &event[k]
			midiEvent := convertToMIDIEvent(e.Data)
			if midiEvent != nil {
				//log.Info.Printf("%s", midiEvent.String())
				j.synth.PushEvent(nil, "midi", midiEvent)
			}
		}
	}

	// read from the audio input buffers
	for i := range j.audioIn {
		audioIn := j.audioIn[i].GetBuffer(nframes)
		copy(j.synth.audio[i][:], audioIn)
	}

	j.synth.Loop()

	// write to the audio output buffers
	ofs := j.synth.nIn
	for i := range j.audioOut {
		audioOut := j.audioOut[i].GetBuffer(nframes)
		copy(audioOut, j.synth.audio[ofs+i][:])
	}

	// write MIDI output events
	// TODO

	return 0
}

func (j *Jack) shutdown() {
	log.Info.Printf("")
}

//-----------------------------------------------------------------------------

// Jack contains state for the jack client.
type Jack struct {
	name     string // client name
	synth    *Synth // top-level synth
	client   *jack.Client
	audioOut []*jack.Port // audio output ports
	audioIn  []*jack.Port // audio input ports
	midiOut  []*jack.Port // midi output ports
	midiIn   []*jack.Port // midi input ports
}

// NewJack returns a jack client object.
func NewJack(name string, synth *Synth) (*Jack, error) {
	j := &Jack{
		name:  name,
		synth: synth,
	}

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

	mi := synth.root.Info()

	// audio output ports
	n := mi.Out.numPortsByType(PortTypeAudio)
	ports, err := j.registerPorts(n, "audio_out", jack.DefaultAudio, jack.PortIsOutput)
	if err != nil {
		j.Close()
		return nil, err
	}
	j.audioOut = ports

	// audio input ports
	n = mi.In.numPortsByType(PortTypeAudio)
	ports, err = j.registerPorts(n, "audio_in", jack.DefaultAudio, jack.PortIsInput)
	if err != nil {
		j.Close()
		return nil, err
	}
	j.audioIn = ports

	// MIDI output ports
	n = mi.Out.numPortsByType(PortTypeMIDI)
	ports, err = j.registerPorts(n, "midi_out", jack.DefaultMIDI, jack.PortIsOutput)
	if err != nil {
		j.Close()
		return nil, err
	}
	j.midiOut = ports

	// MIDI input ports
	n = mi.In.numPortsByType(PortTypeMIDI)
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

//-----------------------------------------------------------------------------

// registerPorts registers a set of jack client ports.
func (j *Jack) registerPorts(n int, prefix, portType string, flags uint64) ([]*jack.Port, error) {
	ports := make([]*jack.Port, n)
	for i := range ports {
		pname := fmt.Sprintf("%s_%d", prefix, i)
		p := j.client.PortRegister(pname, portType, flags, 0)
		if p == nil {
			return nil, fmt.Errorf("can't register port %s:%s", j.name, pname)
		}
		log.Info.Printf("registered \"%s\"", p.Name())
		ports[i] = p
	}
	return ports, nil
}

// unregisterPorts unregisters a set of jack clientports.
func (j *Jack) unregisterPorts(ports []*jack.Port) {
	if j.client == nil || ports == nil {
		return
	}
	for _, p := range ports {
		log.Info.Printf("unregistered \"%s\"", p.Name())
		j.client.PortUnregister(p)
	}
}

//-----------------------------------------------------------------------------
