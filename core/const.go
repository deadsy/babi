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
const AudioBufferSize = 128

// SecsPerAudioBuffer is the audio duration for a single audio buffer.
const SecsPerAudioBuffer = float32(AudioBufferSize) / float32(AudioSampleFrequency)

// FullCycle is a full uint32 phase count.
const FullCycle = 1 << 32

// HalfCycle is a half uint32 phase count.
const HalfCycle = 1 << 31

// FrequencyScale scales a frequency value to a uint32 phase step value.
const FrequencyScale = float32(FullCycle) / float32(AudioSampleFrequency)

// PhaseScale scales a phase value to a uint32 phase step value.
const PhaseScale = float32(FullCycle) / Tau

//-----------------------------------------------------------------------------

// SecsPerMin (seconds per minute).
const SecsPerMin = 60.0

//-----------------------------------------------------------------------------

// MinBeatsPerMin for sequencer.
const MinBeatsPerMin = 10.0

// MaxBeatsPerMin for sequencer.
const MaxBeatsPerMin = 300.0

//-----------------------------------------------------------------------------

// Pi (3.14159...)
const Pi = math.Pi

// Tau (2 * Pi).
const Tau = 2 * math.Pi

//-----------------------------------------------------------------------------
