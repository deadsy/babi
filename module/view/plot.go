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

var plotViewInfo = core.ModuleInfo{
	Name: "plotView",
	In: []core.PortInfo{
		{"x", "x-input", core.PortTypeAudio, nil},
		{"y0", "y-input 0", core.PortTypeAudio, nil},
		{"trigger", "trigger", core.PortTypeBool, plotViewTrigger},
	},
	Out: nil,
}

// Info returns the module information.
func (m *plotView) Info() *core.ModuleInfo {
	return &m.info
}

//-----------------------------------------------------------------------------

type plotView struct {
	info        core.ModuleInfo // module info
	cfg         *PlotConfig     // plot configuration
	x           uint64          // current x-value
	samples     int             // number of samples to plot per trigger
	samplesLeft int             // samples left in this trigger
	idx         int             // file index number
	triggered   bool            // are we currently triggered?
	file        *os.File        // output file
	buf         *bufio.Writer   // buffered io to output file
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
		cfg.Title = "Plot"
	}
	if cfg.X == "" {
		cfg.X = "time"
	}
	if cfg.Y0 == "" {
		cfg.Y0 = "Y0"
	}
	// set the sampling duration
	var samples int
	if cfg.Duration <= 0 {
		// get N buffers of samples
		samples = 4 * core.AudioBufferSize
	} else {
		samples = core.Max(16, int(cfg.Duration/core.AudioSamplePeriod))
	}

	m := &plotView{
		info:    plotViewInfo,
		cfg:     cfg,
		samples: samples,
	}
	return s.Register(m)
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
	trigger := e.GetEventBool().Val
	if !trigger {
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
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *plotView) Process(buf ...*core.Buf) bool {

	if m.triggered {
		x := buf[0]
		y0 := buf[1]
		// how many samples should we plot?
		n := core.Min(m.samplesLeft, core.AudioBufferSize)
		// plot x
		if x != nil {
			m.appendData(m.cfg.X, x[:n])
		} else {
			// no x data - use the internal timebase
			time := make([]float32, n)
			base := float32(m.x) * core.AudioSamplePeriod
			for i := range time {
				time[i] = base
				base += core.AudioSamplePeriod
			}
			m.appendData(m.cfg.X, time)
		}

		// plot y
		if y0 != nil {
			m.appendData(m.cfg.Y0, y0[:n])
		}
		m.samplesLeft -= n
		// are we done?
		if m.samplesLeft == 0 {
			m.closePlot()
			m.triggered = false
		}
	}

	// increment the internal time base
	m.x += core.AudioBufferSize
	return false
}

//-----------------------------------------------------------------------------

// plotName returns the name of the current plot file.
func (m *plotView) plotName() string {
	return fmt.Sprintf("%s%d", m.cfg.Name, m.idx)
}

// open a plot file.
func (m *plotView) openPlot() error {
	log.Info.Printf("open %s.py", m.plotName())
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
	log.Info.Printf("close %s.py", m.plotName())
	// add footer
	m.buf.WriteString(m.footer() + "\n")
	m.buf.Flush()
	m.file.Close()
	m.idx++
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
	s = append(s, fmt.Sprintf("\t\ttitle='%s',", m.cfg.X))
	s = append(s, "\t),")
	s = append(s, "\tyaxis=dict(")
	s = append(s, fmt.Sprintf("\t\ttitle='%s',", m.cfg.Y0))
	s = append(s, "\t\trangemode='tozero',")
	s = append(s, "\t),")
	s = append(s, ")")

	s = append(s, "figure = plotly.graph_objs.Figure(data=data, layout=layout)")
	s = append(s, fmt.Sprintf("plotly.offline.plot(figure, filename='%s.html')", m.plotName()))

	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------
