package dx

type curveType uint

const (
	curveNegLin curveType = 0
	curveNegExp           = 1
	curvePosExp           = 2
	curvePosLin           = 3
)

type oscModeType uint

const (
	oscModeRatio oscModeType = 0
	oscModeFixed             = 1
)

type opParms struct {
	rate              [4]uint // 0..99
	level             [4]uint // 0..99
	brkPoint          uint    // C3 = $27
	lDepth            uint
	rDepth            uint
	lCurve            curveType
	rCurve            curveType
	kbdRateScaling    uint // 0-7
	ampModSensitivity uint // 0-3
	keyVelSensitivity uint // 0-7
	opOutputLevel     uint // 0-99
	oscMode           oscModeType
	oscFreqCoarse     uint // 0-31
	oscFreqFine       uint // 0-99
	oscDetune         uint // 0-14 0: det=-7
}

// SysEx parses a buffer of DX system exclusive MIDI data.
func SysEx(data []byte) {
}
