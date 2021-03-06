//-----------------------------------------------------------------------------
/*

Utility Testing

*/
//-----------------------------------------------------------------------------

package core

import (
	"testing"
)

//-----------------------------------------------------------------------------

func Test_SignExtend(t *testing.T) {
	tests := []struct {
		val    int
		n      uint
		result int
	}{
		{0, 0, 0},
		{1, 0, 0},
		{7, 2, -1},
		{7, 3, -1},
		{3, 3, 3},
		{3, 2, -1},
		{256, 8, 0},
		{256, 9, -256},
		{256, 10, 256},
	}
	for _, v := range tests {
		x := SignExtend(v.val, v.n)
		if v.result != x {
			t.Logf("failed SignExtend(%d, %d) expected %d, got %d", v.val, v.n, v.result, x)
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------
