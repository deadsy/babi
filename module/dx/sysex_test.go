//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package dx

import (
	"fmt"
	"io/ioutil"
	"testing"
)

//-----------------------------------------------------------------------------

func Test_Parse(t *testing.T) {

	tests := []struct {
		path string
	}{
		{"./test/rom1a.syx"},
		{"./test/021.syx"},
		{"./test/060.syx"},
	}

	for _, v := range tests {
		t.Logf("read file %s", v.path)
		buf, err := ioutil.ReadFile(v.path)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		n, err := DecodeSysex(buf)

		if err != nil {
			t.Error(err)
		}

		fmt.Printf("%d bytes\n", n)

	}

}

//-----------------------------------------------------------------------------
