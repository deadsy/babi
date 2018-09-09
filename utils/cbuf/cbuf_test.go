//-----------------------------------------------------------------------------
/*

Circular Buffer Testing

*/
//-----------------------------------------------------------------------------

package cbuf

import (
	"testing"
)

//-----------------------------------------------------------------------------

// Test0 tests basic operations
func Test0(t *testing.T) {

	c := NewCircularBuffer(5)
	x0 := "test string"
	if !c.Empty() {
		t.Error("FAIL")
	}
	// read should fail
	out, err := c.Read()
	if err == nil {
		t.Error("FAIL")
	}
	// write
	err = c.Write(x0)
	if err != nil {
		t.Error("FAIL")
	}
	if c.Empty() {
		t.Error("FAIL")
	}
	// read
	out, err = c.Read()
	if err != nil {
		t.Error("FAIL")
	}
	// check
	if out != x0 {
		t.Error("FAIL")
	}
	if !c.Empty() {
		t.Error("FAIL")
	}
}

//-----------------------------------------------------------------------------

// Test1 tests writes until full
func Test1(t *testing.T) {
	n := 1025
	c := NewCircularBuffer(n)
	// writing 0 thru n-2 should be OK
	for i := 0; i < n-1; i++ {
		err := c.Write(i)
		if err != nil {
			t.Error("FAIL")
		}
	}
	// subsequent writes should fail
	for i := 0; i < n; i++ {
		err := c.Write(i)
		if err == nil {
			t.Error("FAIL")
		}
	}
	// read the values back and check
	for i := 0; i < n-1; i++ {
		x, err := c.Read()
		if err != nil || x != i {
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------
