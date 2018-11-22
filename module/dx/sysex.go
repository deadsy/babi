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
	"unsafe"

	"github.com/deadsy/babi/core"
)

//-----------------------------------------------------------------------------

const midiStatusSysexStart = 0xf0
const midiStatusSysexEnd = 0xf7
const midiIDYamaha = 0x43

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

func (o *opData) convert(idx int) *opConfig {
	cfg := &opConfig{
		idx: idx,
	}

	if o.breakPoint < 3 {
		cfg.breakPoint = core.MidiNote(0)
	} else {
		cfg.breakPoint = core.MidiNote(o.breakPoint - 3)
	}

	for i := 0; i < 4; i++ {
		cfg.env.rate[i] = int(o.rate[i])
		cfg.env.level[i] = int(o.level[i])
	}

	cfg.oscMode = oscModeType(o.x3 & 1)
	cfg.freqCoarse = int((o.x3 >> 1) & 31)
	cfg.freqFine = int(o.freqFine)
	cfg.keyRateScale = int(o.x1 & 7)
	cfg.detune = int((o.x1>>3)&15) - 7
	cfg.outputLevel = int(o.outputLevel)
	cfg.velocitySensitivity = int((o.x2 >> 2) & 7)
	cfg.amSensitivity = int(o.x2 & 3)
	cfg.leftDepth = int(o.leftDepth)
	cfg.rightDepth = int(o.rightDepth)
	cfg.leftCurve = curveType(o.x0 & 3)
	cfg.rightCurve = curveType((o.x0 >> 2) & 3)

	return cfg
}

//-----------------------------------------------------------------------------

type voice128Data struct {
	op               [6]opData // 0: 6..1
	pRate            [4]byte   // 102:
	pLevel           [4]byte   // 106:
	algorithm        byte      // 110: 0..31
	keySyncFeedback  byte      // 111: 0..7
	lfoSpeed         byte      // 112:
	lfoDelay         byte      // 113:
	lfoPhaseModDepth byte      // 114: 0..99
	lfoAmpModDepth   byte      // 115: 0..99
	x2               byte      // 116: 0 ppp www s, pms 0..7, wave 0..7, sync 0..1
	transpose        byte      // 117:
	name             [10]byte  // 118:
}

func (v *voice128Data) convert() *voiceConfig {
	cfg := &voiceConfig{}

	for i := range v.op {
		idx := 5 - i
		cfg.op[idx] = v.op[i].convert(idx)
	}

	for i := 0; i < 4; i++ {
		cfg.env.rate[i] = int(v.pRate[i])
		cfg.env.level[i] = int(v.pLevel[i])
	}

	cfg.algorithm = int(v.algorithm)
	cfg.feedback = int(v.keySyncFeedback & 7)
	cfg.lfo.wave = lfoWaveType((v.x2 >> 1) & 7)
	cfg.lfo.speed = int(v.lfoSpeed)
	cfg.lfo.delay = int(v.lfoDelay)
	cfg.lfo.pmDepth = int(v.lfoPhaseModDepth)
	cfg.lfo.amDepth = int(v.lfoAmpModDepth)
	cfg.lfo.pms = int((v.x2 >> 4) & 7)
	cfg.lfo.sync = v.x2&1 != 0
	cfg.transpose = core.MidiNote(v.transpose + 12)
	cfg.name = string(v.name[:])

	return cfg
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
		return 0, fmt.Errorf("bad checksum: is 0x%02x, should be 0x%02x", csum, buf[n-1])
	}

	voices := (*voices32)(unsafe.Pointer(&buf[0]))

	for i := range voices {
		cfg := voices[i].convert()
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
