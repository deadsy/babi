//-----------------------------------------------------------------------------
/*

Buffer Testing

*/
//-----------------------------------------------------------------------------

package core

import (
	"testing"
)

//-----------------------------------------------------------------------------

func Test_Mul_SS(t *testing.T) {

	a := NewBuf(3)
	b := NewBuf(4)
	c := NewBuf(12)

	a.Mul(b)

	if !a.Equals(c) {
		t.Error("FAIL")
	}

}

//-----------------------------------------------------------------------------
