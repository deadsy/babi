//-----------------------------------------------------------------------------
/*

Constants

*/
//-----------------------------------------------------------------------------

package core

import "math"

//-----------------------------------------------------------------------------

// Audio Sampling Frequency (Hz).
const AudioSampleFrequency = 48000

// Audio Sample Period (seconds).
const AudioSamplePeriod = 1 / AudioSampleFrequency

// Number of float32 samples per audio buffer.
const AudioBufferSize = 64

// Seconds per audio buffer.
const SecsPerAudioBuffer = AudioBufferSize / AudioSampleFrequency

//-----------------------------------------------------------------------------

// Seconds per minute.
const SecsPerMin = 60

//-----------------------------------------------------------------------------

const MinBeatsPerMin = 10
const MaxBeatsPerMin = 300

//-----------------------------------------------------------------------------

const PI = math.Pi
const TAU = 2 * math.Pi

//-----------------------------------------------------------------------------
