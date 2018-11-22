//-----------------------------------------------------------------------------
/*


 */
//-----------------------------------------------------------------------------

package dx

import (
	"math"

	"github.com/deadsy/babi/core"
)

//-----------------------------------------------------------------------------

func midiNoteToLogFreq(note int) int {
	const base = 50857777 // (1 << 24) * (log(440) / log(2) - 69/12)
	const step = (1 << 24) / 12
	return base + step*note
}

var coarsemul = []int{
	-16777216, 0, 16777216, 26591258, 33554432, 38955489, 43368474, 47099600,
	50331648, 53182516, 55732705, 58039632, 60145690, 62083076, 63876816,
	65546747, 67108864, 68576247, 69959732, 71268397, 72509921, 73690858,
	74816848, 75892776, 76922906, 77910978, 78860292, 79773775, 80654032,
	81503396, 82323963, 83117622,
}

func oscFreq(note, mode, coarse, fine, detune int) int {
	// TODO: pitch randomization
	var logfreq int
	if mode == 0 {
		logfreq = midiNoteToLogFreq(note)
		// could use more precision, closer enough for now. those numbers comes from my DX7
		detuneRatio := 0.0209 * math.Exp(-0.396*((float64(logfreq))/(1<<24))) / 7
		logfreq += int(detuneRatio * float64(logfreq) * float64(detune-7))
		logfreq += coarsemul[coarse&31]
		if fine != 0 {
			// (1 << 24) / log(2)
			logfreq += int(math.Floor(24204406.323123*math.Log(1+0.01*float64(fine)) + 0.5))
		}
		// This was measured at 7.213Hz per count at 9600Hz, but the exact
		// value is somewhat dependent on midinote. Close enough for now.
		//logfreq += 12606 * (detune -7);
	} else {
		// ((1 << 24) * log(10) / log(2) * .01) << 3
		logfreq = (4458616 * ((coarse&3)*100 + fine)) >> 3
		if detune > 7 {
			logfreq += 13457 * (detune - 7)
		}
	}
	return logfreq
}

//-----------------------------------------------------------------------------

// freq returns the fixed/ratio frequency for the operator.
func (o *opConfig) freq() float32 {
	var f float32
	switch o.oscMode {
	case oscModeRatio:
		f = float32(o.freqCoarse)
		if f == 0 {
			f = 0.5
		}
		f *= (1.0 + (float32(o.freqFine) * 0.01))
	case oscModeFixed:
		logfreq := (4458616 * ((o.freqCoarse&3)*100 + o.freqFine)) >> 3
		lf := float32(logfreq) * (1.0 / float32(1<<24))
		f = core.Pow2(lf)
	}
	return f
}

//-----------------------------------------------------------------------------
