//-----------------------------------------------------------------------------
/*

Parse DX7 System Exclusive Buffers

https://github.com/asb2m10/dexed/blob/master/Documentation/sysex-format.txt

*/
//-----------------------------------------------------------------------------

package dx

import (
	"errors"
	"fmt"
	"strings"
	"unsafe"

	"github.com/deadsy/babi/core"
)

//-----------------------------------------------------------------------------

const midiStatusSysexStart = 0xf0
const midiStatusSysexEnd = 0xf7
const midiIDYamaha = 0x43

//-----------------------------------------------------------------------------

type curveType int

const (
	curveNegLin curveType = 0
	curveNegExp           = 1
	curvePosExp           = 2
	curvePosLin           = 3
)

func (t curveType) String() string {
	names := []string{"-lin", "-exp", "+exp", "+lin"}
	return names[t]
}

type oscModeType int

const (
	oscModeRatio oscModeType = 0
	oscModeFixed             = 1
)

func (t oscModeType) String() string {
	names := []string{"ratio", "fixed"}
	return names[t]
}

//-----------------------------------------------------------------------------

type opConfig struct {
	idx                 int // operator index 0..5
	outputLevel         int // 0..99
	velocitySensitivity int // 0..7
	amSensitivity       int // 0..3
	oscMode             oscModeType
	freqCoarse          int    // 0..31
	freqFine            int    // 0..99
	detune              int    // -7..7
	rate                [4]int // 0..99
	level               [4]int // 0..99
	breakPoint          core.MidiNote
	keyRateScale        int // 0..7
	leftCurve           curveType
	leftDepth           int // 0..99
	rightCurve          curveType
	rightDepth          int // 0..99
}

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

// row returns a set of row values for this operator.
// needs to match the column order in (v *voiceConfig) String().
func (o *opConfig) row() []string {

	row := make([]string, 0, 16)
	row = append(row, fmt.Sprintf("op%d", o.idx+1))
	row = append(row, fmt.Sprintf("%s", o.oscMode))
	row = append(row, fmt.Sprintf("%.3f", o.freq()))
	//row = append(row, fmt.Sprintf("%d", o.freqCoarse))
	//row = append(row, fmt.Sprintf("%d", o.freqFine))
	row = append(row, fmt.Sprintf("%d", o.detune))
	row = append(row, fmt.Sprintf("[%d %d %d %d]", o.rate[0], o.rate[1], o.rate[2], o.rate[3]))
	row = append(row, fmt.Sprintf("[%d %d %d %d]", o.level[0], o.level[1], o.level[2], o.level[3]))
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
	// pm sensitivity ?
	row = append(row, fmt.Sprintf("%d", o.amSensitivity))

	return row
}

//-----------------------------------------------------------------------------

type voiceConfig struct {
	name      string
	algorithm int
	op        [6]*opConfig
}

func (v *voiceConfig) String() string {
	var s []string
	s = append(s, v.name)
	s = append(s, fmt.Sprintf("algorithm %d", v.algorithm+1))

	rows := make([][]string, len(v.op)+1)
	hdr := []string{
		"",
		"oscMode",
		"freq",
		//"fCoarse",
		//"fFine",
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
	rows[0] = hdr
	for i := range v.op {
		rows[i+1] = v.op[i].row()
	}

	s = append(s, core.TableString(rows, nil, 1))

	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------

type opData struct {
	rate        [4]byte // 0: 0..99
	level       [4]byte // 4: 0..99
	breakPoint  byte    // 8: C3 = $27
	leftDepth   byte    // 9: 0..99
	rightDepth  byte    // 10: 0..99
	x0          byte    // 11: 0000 rr ll, right curve 0..3, left curve 0..3
	x1          byte    // 12: 0 dddd sss, detune 0..14, rate scale 0..7
	x2          byte    // 13: ... kkk aa, key velocity sensitivity 0..7, amp mod sensitivity 0..3
	outputLevel byte    // 14: 0..99
	x3          byte    // 15: 00 fffff m, frequency coarse 0..31, osc mode 0..1
	freqFine    byte    // 16: 0..99
}

func (o *opData) convert(idx int) (*opConfig, error) {
	cfg := &opConfig{
		idx: idx,
	}

	// output level
	if o.outputLevel > 99 {
		return nil, fmt.Errorf("output level is out of range: %d > 99", o.outputLevel)
	}
	cfg.outputLevel = int(o.outputLevel)

	// velocity sensitivity
	cfg.velocitySensitivity = int((o.x2 >> 2) & 7)

	// lfo sensitivity (format?)
	cfg.amSensitivity = int(o.x2 & 3)

	// break point
	if o.breakPoint < 3 {
		cfg.breakPoint = core.MidiNote(0)
	} else {
		cfg.breakPoint = core.MidiNote(o.breakPoint - 3)
	}

	// oscillator mode
	cfg.oscMode = oscModeType(o.x3 & 1)

	// frequency coarse
	cfg.freqCoarse = int((o.x3 >> 1) & 31)

	// frequency fine
	if o.freqFine > 99 {
		return nil, fmt.Errorf("frequency fine is out of range: %d > 99", o.freqFine)
	}
	cfg.freqFine = int(o.freqFine)

	// key rate scale
	cfg.keyRateScale = int(o.x1 & 7)

	// detune
	detune := int((o.x1>>3)&15) - 7
	if detune > 7 {
		return nil, fmt.Errorf("detune is out of range: %d > 7", detune)
	}
	cfg.detune = detune

	// envelope rates and levels
	for i := 0; i < 4; i++ {
		if o.rate[i] > 99 {
			return nil, fmt.Errorf("rate is out of range: %d > 99", o.rate[i])
		}
		cfg.rate[i] = int(o.rate[i])
		if o.level[i] > 99 {
			return nil, fmt.Errorf("level is out of range: %d > 99", o.rate[i])
		}
		cfg.level[i] = int(o.level[i])
	}

	// left/right depth
	if o.leftDepth > 99 {
		return nil, fmt.Errorf("left depth is out of range: %d > 99", o.leftDepth)
	}
	cfg.leftDepth = int(o.leftDepth)

	if o.rightDepth > 99 {
		return nil, fmt.Errorf("right depth is out of range: %d > 99", o.rightDepth)
	}
	cfg.rightDepth = int(o.rightDepth)

	// left/right curve
	cfg.leftCurve = curveType(o.x0 & 3)
	cfg.rightCurve = curveType((o.x0 >> 2) & 3)

	return cfg, nil
}

//-----------------------------------------------------------------------------

type voice128Data struct {
	op              [6]opData // 0: 6..1
	pRate           [4]byte   // 102:
	pLevel          [4]byte   // 106:
	algorithm       byte      // 110: 0..31
	keySyncFeedback byte      // 111:
	lfoSpeed        byte      // 112:
	lfoDelay        byte      // 113:
	x0              byte      // 114: LPMD             LF PT MOD DEP 0-99
	x1              byte      // 115: LAMD             LF AM MOD DEP 0-99
	x2              byte      // 116: |  LPMS |      LFW      |LKS| LF PT MOD SNS 0-7   WAVE 0-5,  SYNC 0-1
	transpose       byte      // 117:
	name            [10]byte  // 118:
}

func (v *voice128Data) convert() (*voiceConfig, error) {
	cfg := &voiceConfig{}

	// name
	cfg.name = string(v.name[:])

	// algorithm
	if v.algorithm > 31 {
		return nil, fmt.Errorf("algorithm is out of range: %d > 31", v.algorithm)
	}
	cfg.algorithm = int(v.algorithm)

	// operators
	for i := range v.op {
		idx := 5 - i
		x, err := v.op[i].convert(idx)
		if err != nil {
			return nil, err
		}
		cfg.op[idx] = x
	}

	return cfg, nil
}

//-----------------------------------------------------------------------------

type voice155Data struct {
	data [155]byte
}

//-----------------------------------------------------------------------------

type voices32 [32]voice128Data

type voicesHdr struct {
	formatNum byte
	countMSB  byte
	countLSB  byte
}

type sysexHdr struct {
	start     byte
	manufID   byte
	subStatus byte
}

//-----------------------------------------------------------------------------

func checksum(buf []byte) byte {
	var csum byte
	for _, c := range buf {
		csum += c
	}
	return -csum & 0x7f
}

func decode32Voice(buf []byte) (int, error) {
	// should have 32 x 128 byte voice records
	n := int(unsafe.Sizeof(voices32{})) + 1
	if len(buf) < n {
		return 0, fmt.Errorf("bad voice data size: is %d, should be %d", len(buf), n)
	}
	// checksum
	csum := checksum(buf[:n-1])
	if csum != buf[n-1] {
		return 0, fmt.Errorf("bad checksum: is 0x%02x, ahould be 0x%02x", csum, buf[n-1])
	}

	voices := (*voices32)(unsafe.Pointer(&buf[0]))

	for i := range voices {
		cfg, err := voices[i].convert()
		if err != nil {
			return 0, err
		}
		fmt.Printf("%s\n", cfg.String())
	}

	return n, nil
}

func decode1Voice(buf []byte) (int, error) {
	// should have a single voice record
	n := int(unsafe.Sizeof(voice155Data{})) + 1
	if len(buf) < n {
		return 0, fmt.Errorf("bad voice data size: is %d, should be %d", len(buf), n)
	}
	csum := checksum(buf[:n-1])
	if csum != buf[n-1] {
		return 0, fmt.Errorf("bad checksum: is 0x%02x, ahould be 0x%02x", csum, buf[n-1])
	}

	return n, nil
}

func decodeVoice(buf []byte) (int, error) {
	ofs := 0
	n := int(unsafe.Sizeof(voicesHdr{}))
	if len(buf) < n {
		return 0, errors.New("voice sysex header is too short")
	}
	hdr := (*voicesHdr)(unsafe.Pointer(&buf[0]))
	ofs += n

	count := (int(hdr.countMSB) << 7) + int(hdr.countLSB)

	switch hdr.formatNum {
	case 9:
		if count != 4096 {
			return 0, fmt.Errorf("bad voice data count: is %d, should be 4096", count)
		}
		n, err := decode32Voice(buf[ofs:])
		if err != nil {
			return 0, err
		}
		ofs += n
	case 0:
		if count != 155 {
			return 0, fmt.Errorf("bad voice data count: is %d, should be 155", count)
		}
		n, err := decode1Voice(buf[ofs:])
		if err != nil {
			return 0, err
		}
		ofs += n
	default:
		return 0, fmt.Errorf("unknown format number: 0x%02x", hdr.formatNum)
	}

	return ofs, nil
}

func decodeParameterChange(buf []byte) (int, error) {
	panic("todo")
	return len(buf), nil
}

// DecodeSysex parses a buffer of DX system exclusive MIDI data.
func DecodeSysex(buf []byte) (int, error) {
	ofs := 0
	n := int(unsafe.Sizeof(sysexHdr{}))
	if len(buf) < n {
		return 0, errors.New("sysex is too short")
	}
	hdr := (*sysexHdr)(unsafe.Pointer(&buf[ofs]))
	ofs += n

	if hdr.start != midiStatusSysexStart {
		return 0, errors.New("bad sysex start byte")
	}

	if hdr.manufID != midiIDYamaha {
		return 0, fmt.Errorf("bad manufacturer id: 0x%02x", hdr.manufID)
	}

	switch hdr.subStatus {
	case 0:
		n, err := decodeVoice(buf[ofs:])
		if err != nil {
			return 0, err
		}
		ofs += n
	case 0x10:
		n, err := decodeParameterChange(buf[ofs:])
		if err != nil {
			return 0, err
		}
		ofs += n
	default:
		return 0, fmt.Errorf("unknown sub status: 0x%02x", hdr.subStatus)
	}

	if len(buf[ofs:]) < 1 {
		return 0, errors.New("no sysex end byte")
	}
	if buf[ofs] != midiStatusSysexEnd {
		return 0, errors.New("bad sysex end byte")
	}
	ofs++

	return ofs, nil
}

//-----------------------------------------------------------------------------
