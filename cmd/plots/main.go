//-----------------------------------------------------------------------------
/*

Graphical Plots of Waveforms

*/
//-----------------------------------------------------------------------------

package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/osc"
	"github.com/deadsy/babi/module/view"
)

//-----------------------------------------------------------------------------

type trace struct {
	file *os.File
	buf  *bufio.Writer
}

func openTrace(name string) (*trace, error) {
	file, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return &trace{
		file: file,
		buf:  bufio.NewWriter(file),
	}, nil
}

func (t *trace) setTitle(title string) {
	t.buf.WriteString(fmt.Sprintf("title = \"%s\"\n", title))
}

func (t *trace) newVariable(name string) {
	t.buf.WriteString(fmt.Sprintf("%s = []\n", name))
}

func (t *trace) appendData(name string, buf *core.Buf) {
	t.buf.WriteString(fmt.Sprintf("%s.extend([\n", name))
	for i := range buf {
		t.buf.WriteString(fmt.Sprintf("%f,\n", buf[i]))
	}
	t.buf.WriteString(fmt.Sprintf("])\n"))
}

func (t *trace) closeTrace() {
	t.buf.Flush()
	t.file.Close()
}

//-----------------------------------------------------------------------------

func main() {

	name := "babi.py"
	f, err := openTrace(name)
	if err != nil {
		fmt.Printf("unable to create %s\n", name)
		os.Exit(1)
	}

	f.setTitle("220 Hz Sine Wave")
	f.newVariable("time")
	f.newVariable("amplitude")

	t := view.NewTime(nil)
	s := osc.NewSine(nil)
	core.SendEventFloat(s, "frequency", 220.0)

	var x core.Buf
	var y core.Buf

	for i := 0; i < 10; i++ {
		t.Process(&x)
		s.Process(&y)
		f.appendData("time", &x)
		f.appendData("amplitude", &y)
	}

	f.closeTrace()
	os.Exit(0)
}

//-----------------------------------------------------------------------------
