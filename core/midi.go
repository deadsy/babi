//-----------------------------------------------------------------------------
/*

MIDI Functions

*/
//-----------------------------------------------------------------------------

package core

import (
	"fmt"

	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

// Channel Messages
const midiStatusNoteOff = 8 << 4
const midiStatusNoteOn = 9 << 4
const midiStatusPolyphonicAftertouch = 10 << 4
const midiStatusControlChange = 11 << 4
const midiStatusProgramChange = 12 << 4
const midiStatusChannelAftertouch = 13 << 4
const midiStatusPitchWheel = 14 << 4

// System Common Messages
const midiStatusSysexStart = 0xf0
const midiStatusQuarterFrame = 0xf1
const midiStatusSongPointer = 0xf2
const midiStatusSongSelect = 0xf3
const midiStatusTuneRequest = 0xf6
const midiStatusSysexEnd = 0xf7

// System Realtime Messages
const midiStatusTimingClock = 0xf8
const midiStatusStart = 0xfa
const midiStatusContinue = 0xfb
const midiStatusStop = 0xfc
const midiStatusActiveSensing = 0xfe
const midiStatusReset = 0xff

// delimiters
const midiStatusCommon = 0xf0
const midiStatusRealtime = 0xf8

// convertToMIDIEvent converts a midi data buffer into a MIDI event.
func convertToMIDIEvent(data []byte) *Event {
	if len(data) == 0 {
		return nil
	}
	status := data[0]
	if status < midiStatusCommon {
		// channel message
		switch status & 0xf0 {
		case midiStatusNoteOff,
			midiStatusNoteOn,
			midiStatusPolyphonicAftertouch,
			midiStatusControlChange,
			midiStatusPitchWheel:
			et := EventTypeMIDI(status & 0xf0)
			if len(data) == 3 {
				return NewEventMIDI(et, data[0], data[1], data[2])
			}
			log.Info.Printf("%s len(data) != 3", midiEventType2String[et])
		case midiStatusProgramChange,
			midiStatusChannelAftertouch:
			et := EventTypeMIDI(status & 0xf0)
			if len(data) == 2 {
				return NewEventMIDI(EventTypeMIDI(status&0xf0), data[0], data[1], 0)
			}
			log.Info.Printf("%s: len(data) != 2", midiEventType2String[et])
		default:
			log.Info.Printf("unhandled channel msg %02x", status)
		}
	} else if status < midiStatusRealtime {
		// system common message
		switch status {
		case midiStatusSysexStart:
		case midiStatusQuarterFrame:
		case midiStatusSongPointer:
		case midiStatusSongSelect:
		case midiStatusTuneRequest:
		case midiStatusSysexEnd:
		default:
			log.Info.Printf("unhandled system commmon msg %02x", status)
		}
	} else {
		// system real time message
		switch status {
		case midiStatusTimingClock:
		case midiStatusStart:
		case midiStatusContinue:
		case midiStatusStop:
		case midiStatusActiveSensing:
		case midiStatusReset:
		default:
			log.Info.Printf("unhandled system realtime msg %02x", status)
			break
		}
	}
	return nil
}

//-----------------------------------------------------------------------------

const notesInOctave = 12

var sharpNotes = [notesInOctave]string{
	"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B",
}

var flatNotes = [notesInOctave]string{
	"C", "Db", "D", "Eb", "E", "F", "Gb", "G", "Ab", "A", "Bb", "B",
}

// MidiNote 0..127
type MidiNote byte

// Octave returns the MIDI octave of the MIDI note.
func (n MidiNote) Octave() int {
	return int(n) / notesInOctave
}

func (n MidiNote) sharpString() string {
	return sharpNotes[n%notesInOctave]
}

func (n MidiNote) flatString() string {
	return flatNotes[n%notesInOctave]
}

func (n MidiNote) String() string {
	return n.sharpString() + fmt.Sprintf("%d", n.Octave())
}

//-----------------------------------------------------------------------------

// MIDIMap maps a 0..127 midi control value from a..b
func MIDIMap(val uint8, a, b float32) float32 {
	return a + ((b-a)/127.0)*float32(val&0x7f)
}

// MIDIPitchBend maps a pitch bend value onto a MIDI note offset.
func MIDIPitchBend(val uint16) float32 {
	// 0..8192..16383 maps to -/+ 2 semitones
	return float32(val-8192) * (2.0 / 8192.0)
}

// MIDIToFrequency converts a MIDI note to a frequency value (Hz).
// The note value is a float for pitch bending, tuning, etc.
func MIDIToFrequency(note float32) float32 {
	return 440.0 * Pow2((note-69.0)*(1.0/12.0))
}

//-----------------------------------------------------------------------------
