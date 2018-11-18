//-----------------------------------------------------------------------------
/*

Graphical Plots of Waveforms

Produces python code viewable using the plot.ly library.

*/
//-----------------------------------------------------------------------------

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/dx"
	"github.com/deadsy/babi/module/osc"
	"github.com/deadsy/babi/module/view"
)

//-----------------------------------------------------------------------------

type plot struct {
	name  string        // name of output file py/html
	title string        // title of plot
	file  *os.File      // output file
	buf   *bufio.Writer // buffered io to output file
}

// openPlot creates and opens a plot file.
func openPlot(name string) (*plot, error) {
	file, err := os.Create(name + ".py")
	if err != nil {
		return nil, err
	}

	p := &plot{
		name: name,
		file: file,
		buf:  bufio.NewWriter(file),
	}

	// add header
	p.buf.WriteString(p.header() + "\n")
	return p, nil
}

// header returns the python plot file header.
func (p *plot) header() string {
	var s []string
	s = append(s, "#!/usr/bin/env python3")
	s = append(s, "import plotly")
	return strings.Join(s, "\n")
}

// footer returns the python plot file footer.
func (p *plot) footer() string {
	var s []string

	s = append(s, "data = [")
	s = append(s, "\tplotly.graph_objs.Scatter(")
	s = append(s, "\t\tx=time,")
	s = append(s, "\t\ty=amplitude,")
	s = append(s, "\t\tmode = 'lines',")
	s = append(s, "\t),")
	s = append(s, "]")

	s = append(s, "layout = plotly.graph_objs.Layout(")
	s = append(s, fmt.Sprintf("\ttitle='%s',", p.title))
	s = append(s, "\txaxis=dict(")
	s = append(s, "\t\ttitle='time',")
	s = append(s, "\t),")
	s = append(s, "\tyaxis=dict(")
	s = append(s, "\t\ttitle='amplitude',")
	s = append(s, "\t\trangemode='tozero',")
	s = append(s, "\t),")
	s = append(s, ")")

	s = append(s, "figure = plotly.graph_objs.Figure(data=data, layout=layout)")
	s = append(s, fmt.Sprintf("plotly.offline.plot(figure, filename='%s.html')", p.name))

	return strings.Join(s, "\n")
}

// setTitle sets the title of the plot file.
func (p *plot) setTitle(title string) {
	p.title = title
}

// newVariable adds a new variable to the plot file.
func (p *plot) newVariable(name string) {
	p.buf.WriteString(name + " = []\n")
}

// appendData appends named data to the plot file.
func (p *plot) appendData(name string, buf *core.Buf) {
	p.buf.WriteString(name + ".extend([\n")
	for i := range buf {
		p.buf.WriteString(fmt.Sprintf("%f,", buf[i]))
		if i&15 == 15 {
			p.buf.WriteString("\n")
		}
	}
	p.buf.WriteString("])\n")
}

// close the plot file.
func (p *plot) close() {
	// add footer
	p.buf.WriteString(p.footer() + "\n")
	p.buf.Flush()
	p.file.Close()
}

//-----------------------------------------------------------------------------

func envDx() {

	name := "babi"
	p, err := openPlot(name)
	if err != nil {
		fmt.Printf("unable to create %s\n", name)
		os.Exit(1)
	}
	defer p.close()

	p.setTitle(fmt.Sprintf("DX7 Envelope"))
	p.newVariable("time")
	p.newVariable("amplitude")

	levels := &[4]int{99, 80, 99, 0}
	rates := &[4]int{80, 80, 70, 80}

	t := view.NewTime(nil)
	e := dx.NewEnv(nil, levels, rates)
	core.SendEventFloat(e, "gate", 1.0)

	var x core.Buf
	var y core.Buf

	for i := 0; i < 12; i++ {
		t.Process(&x)
		e.Process(&y)
		p.appendData("time", &x)
		p.appendData("amplitude", &y)
	}

	core.SendEventFloat(e, "gate", 0.0)

	for i := 0; i < 4; i++ {
		t.Process(&x)
		e.Process(&y)
		p.appendData("time", &x)
		p.appendData("amplitude", &y)
	}
}

//-----------------------------------------------------------------------------

func goom() {

	name := "babi"
	p, err := openPlot(name)
	if err != nil {
		fmt.Printf("unable to create %s\n", name)
		os.Exit(1)
	}
	defer p.close()

	freq := float32(110.0)
	p.setTitle(fmt.Sprintf("%.1f Hz Goom Wave", freq))
	p.newVariable("time")
	p.newVariable("amplitude")

	t := view.NewTime(nil)
	s := osc.NewGoom(nil)
	core.SendEventFloat(s, "frequency", freq)
	core.SendEventFloat(s, "duty", 0.3)
	core.SendEventFloat(s, "slope", 1.0)

	var x core.Buf
	var y core.Buf

	for i := 0; i < 10; i++ {
		t.Process(&x)
		s.Process(&y)
		p.appendData("time", &x)
		p.appendData("amplitude", &y)
	}
}

//-----------------------------------------------------------------------------

func main() {
	//goom()
	envDx()
	os.Exit(0)
}

//-----------------------------------------------------------------------------
