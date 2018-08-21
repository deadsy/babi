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
// events

type EventType uint

const (
	Event_Null EventType = iota
	Event_MIDI
)

type Event struct {
	Etype EventType   // event type
	Info  interface{} // event information
}

//-----------------------------------------------------------------------------

type Patch interface {
	Process(in, out []*Buf) // run the patch
	Event(e *Event)         // process an event
	Active() bool           // return true if the patch has non-zero output
}

//-----------------------------------------------------------------------------
// voices (active patches)

/*

type Voice struct {
	channel uint  // channel for this voice
	note    uint  // base note for this voice
	patch   Patch // active patch for this voice
}

// NewVoice creates a new active voice.
func (s *Synth) NewVoice(channel, note uint) *Voice {
	return &Voice{
		channel: channel,
		note:    note,
		patch:   s.patch[channel].New(s),
	}
}

// Allocate and assign a voice to a channel.
func (s *Synth) VoiceAlloc(channel, note uint) error {
	// validate the channel
	if channel >= MAX_CHANNELS || s.patch[channel] == nil {
		return fmt.Errorf("no patch defined for channel %d", channel)
	}

	// stop any pre-existing voice in this slot
	v := s.voice[s.voice_idx]
	if v != nil {
		v.patch.Stop()
	}

	// allocate and start a new voice
	v = s.NewVoice(channel, note)
	s.voice[s.voice_idx] = v

	// move to the next voice slot
	s.voice_idx += 1
	if s.voice_idx == MAX_VOICES {
		s.voice_idx = 0
	}

	return nil
}

// Lookup the voice being used for this channel and note.
func (s *Synth) VoiceLookup(channel, note uint) *Voice {
	for i := 0; i < MAX_VOICES; i++ {
		v := s.voice[i]
		if v != nil && v.channel == channel && v.note == note {
			return v
		}
	}
	return nil
}

*/

//-----------------------------------------------------------------------------
// patches

/*

type Patch interface {
	Stop()        // stop the patch
	Process()     // run the patch
	Active() bool // is this patch active?
}



// Add a patch to a channel
func (s *Synth) AddPatch(patch *PatchInfo, channel uint) error {
	if channel >= MAX_CHANNELS {
		return fmt.Errorf("channel %d is out of range", channel)
	}
	if s.patch[channel] != nil {
		return fmt.Errorf("patch %s is already set for channel %d", s.patch[channel].Name, channel)
	}
	s.patch[channel] = patch
	return nil
}

*/

//-----------------------------------------------------------------------------

type Synth struct {
	root  Patch               // root patch
	audio Audio               // audio output device
	out   [AUDIO_CHANNELS]Buf // audio output buffers
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
	for {
		// zero the audio output buffers
		for i := 0; i < AUDIO_CHANNELS; i++ {
			s.out[i].Zero()
		}
		// process the patches
		if s.root != nil && s.root.Active() {
			s.root.Process(nil, nil)
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
