//-----------------------------------------------------------------------------
/*

MIDI Functions

*/
//-----------------------------------------------------------------------------

package core

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
