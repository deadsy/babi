//-----------------------------------------------------------------------------
/*

DX7 Configuration

*/
//-----------------------------------------------------------------------------

package dx

import (
	"fmt"
	"strings"

	"github.com/deadsy/babi/core"
)

//-----------------------------------------------------------------------------
// envelope generator configuration

type envConfig struct {
	rate  [4]int // 0..99
	level [4]int // 0..99
}

func (e *envConfig) rateString() string {
	return fmt.Sprintf("[%d %d %d %d]", e.rate[0], e.rate[1], e.rate[2], e.rate[3])
}

func (e *envConfig) levelString() string {
	return fmt.Sprintf("[%d %d %d %d]", e.level[0], e.level[1], e.level[2], e.level[3])
}

//-----------------------------------------------------------------------------
// low frequency oscillator configuration

type lfoWaveType int

const (
	lfoTriangle      lfoWaveType = 0
	lfoSawDown                   = 1
	lfoSawUp                     = 2
	lfoSquare                    = 3
	lfoSine                      = 4
	lfoSampleAndHold             = 5
)

func (t lfoWaveType) String() string {
	return []string{"tri", "sw-", "sw+", "sqr", "sin", "s&h"}[t]
}

type lfoConfig struct {
	wave    lfoWaveType
	speed   int
	delay   int
	pmDepth int  // phase modulation depth
	amDepth int  // amplitude modulation depth
	pms     int  // pitch modulation sensitivity
	sync    bool // resync lfo with note on event
}

func (l *lfoConfig) String() string {
	rows := make([][]string, 2)
	rows[0] = []string{"lfo", "wave", "speed", "delay", "pmDepth", "amDepth", "pms", "sync"}
	var s []string
	s = append(s, "")
	s = append(s, fmt.Sprintf("%s", l.wave))
	s = append(s, fmt.Sprintf("%d", l.speed))
	s = append(s, fmt.Sprintf("%d", l.delay))
	s = append(s, fmt.Sprintf("%d", l.pmDepth))
	s = append(s, fmt.Sprintf("%d", l.amDepth))
	s = append(s, fmt.Sprintf("%d", l.pms))
	s = append(s, fmt.Sprintf("%s", core.BoolToString(l.sync, []string{"off", "on"})))
	rows[1] = s
	return core.TableString(rows, nil, 1)
}

//-----------------------------------------------------------------------------
// operator configuration

type curveType int

const (
	curveNegLin curveType = 0
	curveNegExp           = 1
	curvePosExp           = 2
	curvePosLin           = 3
)

func (t curveType) String() string {
	return []string{"-lin", "-exp", "+exp", "+lin"}[t]
}

type oscModeType int

const (
	oscModeRatio oscModeType = 0
	oscModeFixed             = 1
)

func (t oscModeType) String() string {
	return []string{"ratio", "fixed"}[t]
}

type opConfig struct {
	idx                 int // operator index 0..5
	outputLevel         int // 0..99
	velocitySensitivity int // 0..7
	amSensitivity       int // 0..3
	oscMode             oscModeType
	freqCoarse          int // 0..31
	freqFine            int // 0..99
	detune              int // -7..7
	env                 envConfig
	breakPoint          core.MidiNote
	keyRateScale        int // 0..7
	leftCurve           curveType
	leftDepth           int // 0..99
	rightCurve          curveType
	rightDepth          int // 0..99
}

var opRowHeader = []string{
	"",
	"oscMode",
	"fCoarse",
	"fFine",
	"detune",
	"rate",
	"level",
	"brkPoint",
	"lCurve",
	"lDepth",
	"rCurve",
	"rDepth",
	"keyRate",
	"outLevel",
	"velSens",
	"amSens",
}

// row returns a set of row values for this operator.
// needs to match the column order in (v *voiceConfig) String().
func (o *opConfig) rowStrings() []string {
	row := make([]string, 0, 16)
	row = append(row, fmt.Sprintf("op%d", o.idx+1))
	row = append(row, fmt.Sprintf("%s", o.oscMode))
	row = append(row, fmt.Sprintf("%d", o.freqCoarse))
	row = append(row, fmt.Sprintf("%d", o.freqFine))
	row = append(row, fmt.Sprintf("%d", o.detune))
	row = append(row, o.env.rateString())
	row = append(row, o.env.levelString())
	if o.breakPoint == core.MidiNote(0) {
		row = append(row, "off")
	} else {
		row = append(row, fmt.Sprintf("%s", o.breakPoint))
	}
	row = append(row, fmt.Sprintf("%s", o.leftCurve))
	row = append(row, fmt.Sprintf("%d", o.leftDepth))
	row = append(row, fmt.Sprintf("%s", o.rightCurve))
	row = append(row, fmt.Sprintf("%d", o.rightDepth))
	row = append(row, fmt.Sprintf("%d", o.keyRateScale))
	row = append(row, fmt.Sprintf("%d", o.outputLevel))
	row = append(row, fmt.Sprintf("%d", o.velocitySensitivity))
	row = append(row, fmt.Sprintf("%d", o.amSensitivity))
	return row
}

//-----------------------------------------------------------------------------
// voice configuration

type voiceConfig struct {
	name      string
	algorithm int
	env       envConfig
	transpose core.MidiNote
	feedback  int
	lfo       lfoConfig
	op        [6]*opConfig
}

func (v *voiceConfig) String() string {
	var s []string
	s = append(s, v.name)
	s = append(s, fmt.Sprintf("transpose %s", v.transpose))
	s = append(s, fmt.Sprintf("rate %s", v.env.rateString()))
	s = append(s, fmt.Sprintf("level %s", v.env.levelString()))
	s = append(s, fmt.Sprintf("algorithm %d", v.algorithm+1))
	s = append(s, fmt.Sprintf("feedback %d", v.feedback))
	s = append(s, fmt.Sprintf("%s", &v.lfo))
	// operators in table form
	rows := make([][]string, len(v.op)+1)
	rows[0] = opRowHeader
	for i := range v.op {
		rows[i+1] = v.op[i].rowStrings()
	}
	s = append(s, core.TableString(rows, nil, 1))
	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------
