//-----------------------------------------------------------------------------
/*

Random Testing

*/
//-----------------------------------------------------------------------------

package core

import (
	"testing"
)

//-----------------------------------------------------------------------------

func Test_Float(t *testing.T) {
	r := NewRand(0)
	for i := 0; i < 100; i++ {
		t.Logf("%f", r.Float())
	}

}

//-----------------------------------------------------------------------------
