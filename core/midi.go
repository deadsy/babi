//-----------------------------------------------------------------------------
/*

MIDI Functions

*/
//-----------------------------------------------------------------------------

package core

//-----------------------------------------------------------------------------

// MIDI_Map maps a 0..127 midi control value from a..b
func MIDI_Map(val uint8, a, b float32) float32 {
	return a + ((b-a)/127.0)*float32(val&0x7f)
}

// MIDI_PitchBend maps a pitch bend value onto a note offset
func MIDI_PitchBend(val uint16) float32 {
	// 0..8192..16383 maps to -/+ 2 semitones
	return float32(val-8192) * (2.0 / 8192.0)
}

// MIDI_ToFrequency does midi note to frequency conversion
// Note: treat the note as a float for pitch bending, tuning, etc.
func MIDI_ToFrequency(note float32) float32 {
	return 440.0 * Pow2((note-69.0)*(1.0/12.0))
}

//-----------------------------------------------------------------------------
