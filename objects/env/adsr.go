//-----------------------------------------------------------------------------
/*

ADSR Envelope Object

*/
//-----------------------------------------------------------------------------

package env

import "github.com/deadsy/babi/core"

//-----------------------------------------------------------------------------

type ADSR struct {
}

func NewADSR() *ADSR {
	return &ADSR{}
}

func (e *ADSR) Active() bool {
	return true
}

func (e *ADSR) Process(out *core.SBuf) {
}

//-----------------------------------------------------------------------------
