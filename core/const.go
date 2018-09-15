//-----------------------------------------------------------------------------
/*

Constants

*/
//-----------------------------------------------------------------------------

package core

import "math"

//-----------------------------------------------------------------------------

// AudioSampleFrequency is the sample frequency for audio (Hz).
const AudioSampleFrequency = 48000

// AudioSamplePeriod is the sample period for audio (seconds).
const AudioSamplePeriod = 1.0 / float32(AudioSampleFrequency)

// AudioBufferSize is the number of float32 samples per audio buffer.
const AudioBufferSize = 64

// SecsPerAudioBuffer is the audio duration for a single audio buffer.
const SecsPerAudioBuffer = float32(AudioBufferSize) / float32(AudioSampleFrequency)

//-----------------------------------------------------------------------------

// SecsPerMin
const SecsPerMin = 60.0

//-----------------------------------------------------------------------------

// MinBeatsPerMin for sequencer.
const MinBeatsPerMin = 10.0

// MaxBeatsPerMin for sequencer.
const MaxBeatsPerMin = 300.0

//-----------------------------------------------------------------------------

// Pi
const PI = math.Pi

// Tau (2 * Pi).
const TAU = 2 * math.Pi

//-----------------------------------------------------------------------------
