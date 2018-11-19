//-----------------------------------------------------------------------------
/*

Plotting Module

When this module is triggered it writes the input signal to python code
that uses the plot.ly library to create a plot of the signal.

*/
//-----------------------------------------------------------------------------

package view

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

// Info returns the module information.
func (m *plotView) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "plotView",
		In: []core.PortInfo{
			{"x", "x-input", core.PortTypeAudioBuffer, nil},
			{"y0", "y-input 0", core.PortTypeAudioBuffer, nil},
			{"trigger", "trigger (!= 0)", core.PortTypeInt, plotViewTrigger},
		},
		Out: nil,
	}
}

//-----------------------------------------------------------------------------

type plotView struct {
	synth       *core.Synth   // top-level synth
	cfg         *PlotConfig   // plot configuration
	samples     int           // number of samples to plot per trigger
	samplesLeft int           // samples left in this trigger
	count       int           // how many times have we been triggered?
	triggered   bool          // are we currently triggered?
	file        *os.File      // output file
	buf         *bufio.Writer // buffered io to output file
}

// PlotConfig provides the configuration for a plotting module.
type PlotConfig struct {
	Name     string  // name of output python file
	Title    string  // title of plot
	X        string  // x-axis name
	Y0       string  // y0-axis name
	Duration float32 // sample time in seconds
}

// NewPlot returns a signal plotting module.
func NewPlot(s *core.Synth, cfg *PlotConfig) core.Module {
	log.Info.Printf("")
	// set some defaults
	if cfg.Name == "" {
		cfg.Name = "plot"
	}
	if cfg.Title == "" {
		cfg.Name = "Plot"
	}
	if cfg.X == "" {
		cfg.Name = "X0"
	}
	if cfg.Y0 == "" {
		cfg.Name = "Y0"
	}
	// set the sampling duration
	var samples int
	if cfg.Duration <= 0 {
		// get N buffers of samples
		samples = 4 * core.AudioBufferSize
	} else {
		samples = core.Max(16, int(cfg.Duration/core.AudioSamplePeriod))
	}

	return &plotView{
		synth:   s,
		cfg:     cfg,
		samples: samples,
	}
}

// Child returns the child modules of this module.
func (m *plotView) Child() []core.Module {
	return nil
}

// Stop performs any cleanup of a module.
func (m *plotView) Stop() {
	if m.triggered {
		m.closePlot()
	}
}

//-----------------------------------------------------------------------------
// Port Events

func plotViewTrigger(cm core.Module, e *core.Event) {
	m := cm.(*plotView)
	trigger := e.GetEventInt().Val
	if trigger == 0 {
		return
	}
	if m.triggered {
		log.Info.Printf("already triggered")
		return
	}
	// trigger!
	err := m.openPlot()
	if err != nil {
		log.Info.Printf("can't open plot \"%s\": %s", m.plotName(), err)
		return
	}
	m.newVariable(m.cfg.X)
	m.newVariable(m.cfg.Y0)
	m.triggered = true
	m.samplesLeft = m.samples
	m.count += 1
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *plotView) Process(buf ...*core.Buf) {
	if !m.triggered {
		return
	}
	x := buf[0]
	y0 := buf[1]
	// how many samples should we plot?
	n := core.Min(m.samplesLeft, core.AudioBufferSize)
	// plot x
	if x != nil {
		m.appendData(m.cfg.X, x[:n])
	}
	// plot y
	if y0 != nil {
		m.appendData(m.cfg.Y0, y0[:n])
	}
	m.samplesLeft -= n
	// are we done?
	if m.samplesLeft == 0 {
		m.triggered = false
		m.closePlot()
	}
}

// Active returns true if the module has non-zero output.
func (m *plotView) Active() bool {
	return m.triggered
}

//-----------------------------------------------------------------------------

// plotName returns the name of the current plot file.
func (m *plotView) plotName() string {
	return fmt.Sprintf("%s%d", m.cfg.Name, m.count)
}

// open a plot file.
func (m *plotView) openPlot() error {
	file, err := os.Create(m.plotName() + ".py")
	if err != nil {
		return err
	}
	m.file = file
	m.buf = bufio.NewWriter(file)
	// add header
	m.buf.WriteString(m.header() + "\n")
	return nil
}

// close the plot file.
func (m *plotView) closePlot() {
	// add footer
	m.buf.WriteString(m.footer() + "\n")
	m.buf.Flush()
	m.file.Close()
}

// newVariable adds a new variable to the plot file.
func (m *plotView) newVariable(name string) {
	m.buf.WriteString(name + " = []\n")
}

// appendData appends named data to the plot file.
func (m *plotView) appendData(name string, buf []float32) {
	m.buf.WriteString(name + ".extend([\n")
	for i := range buf {
		m.buf.WriteString(fmt.Sprintf("%f,", buf[i]))
		if i&15 == 15 {
			m.buf.WriteString("\n")
		}
	}
	m.buf.WriteString("])\n")
}

// header returns the python plot file header.
func (m *plotView) header() string {
	var s []string
	s = append(s, "#!/usr/bin/env python3")
	s = append(s, "import plotly")
	return strings.Join(s, "\n")
}

// footer returns the python plot file footer.
func (m *plotView) footer() string {
	var s []string

	s = append(s, "data = [")
	s = append(s, "\tplotly.graph_objs.Scatter(")
	s = append(s, "\t\tx=time,")
	s = append(s, "\t\ty=amplitude,")
	s = append(s, "\t\tmode = 'lines',")
	s = append(s, "\t),")
	s = append(s, "]")

	s = append(s, "layout = plotly.graph_objs.Layout(")
	s = append(s, fmt.Sprintf("\ttitle='%s',", m.cfg.Title))
	s = append(s, "\txaxis=dict(")
	s = append(s, "\t\ttitle='time',")
	s = append(s, "\t),")
	s = append(s, "\tyaxis=dict(")
	s = append(s, "\t\ttitle='amplitude',")
	s = append(s, "\t\trangemode='tozero',")
	s = append(s, "\t),")
	s = append(s, ")")

	s = append(s, "figure = plotly.graph_objs.Figure(data=data, layout=layout)")
	s = append(s, fmt.Sprintf("plotly.offline.plot(figure, filename='%s.html')", m.plotName()))

	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------
