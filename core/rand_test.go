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

func TestFloat32(t *testing.T) {
	buckets := make([]int, 10)
	r := NewRand32(0)
	for i := 0; i < 10000; i++ {
		xf := (r.Float32() + 1.0) * 5.0
		xi := int(xf)
		buckets[xi]++
	}
	for i, v := range buckets {
		t.Logf("%d: %d", i, v)
	}
}

func TestFloat64(t *testing.T) {
	buckets := make([]int, 10)
	r := NewRand64(0)
	for i := 0; i < 10000; i++ {
		xf := (r.Float64() + 1.0) * 5.0
		xi := int(xf)
		buckets[xi]++
	}
	for i, v := range buckets {
		t.Logf("%d: %d", i, v)
	}
}

//-----------------------------------------------------------------------------
