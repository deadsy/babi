//-----------------------------------------------------------------------------
/*

Graphical Plots of Waveforms

Produces python code viewable using the plot.ly library.

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"
	"os"

	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/dx"
	"github.com/deadsy/babi/module/osc"
	"github.com/deadsy/babi/module/view"
)

//-----------------------------------------------------------------------------

func envDx() {
	cfg := &view.PlotConfig{
		Name:     "envDx",
		Title:    fmt.Sprintf("DX7 Envelope"),
		Y0:       "amplitude",
		Duration: 2.0,
	}

	levels := &[4]int{99, 80, 99, 0}
	rates := &[4]int{80, 80, 70, 80}

	s := dx.NewEnv(nil, levels, rates)
	core.SendEventFloat(s, "gate", 1.0)

	p := view.NewPlot(nil, cfg)
	core.SendEventInt(p, "trigger", 1)

	for i := 0; i < 12; i++ {
		var y core.Buf
		s.Process(&y)
		p.Process(nil, &y)
	}

	core.SendEventFloat(s, "gate", 0.0)

	for i := 0; i < 4; i++ {
		var y core.Buf
		s.Process(&y)
		p.Process(nil, &y)
	}

	p.Stop()
}

//-----------------------------------------------------------------------------

func goom() {
	freq := float32(110.0)

	cfg := &view.PlotConfig{
		Name:     "goom",
		Title:    fmt.Sprintf("%.1f Hz Goom Wave", freq),
		Y0:       "amplitude",
		Duration: 2.0,
	}

	s := osc.NewGoom(nil)
	core.SendEventFloat(s, "frequency", freq)
	core.SendEventFloat(s, "duty", 0.3)
	core.SendEventFloat(s, "slope", 1.0)

	p := view.NewPlot(nil, cfg)
	core.SendEventInt(p, "trigger", 1)

	for i := 0; i < 10; i++ {
		var y core.Buf
		s.Process(&y)
		p.Process(nil, &y)
	}

	p.Stop()
}

//-----------------------------------------------------------------------------

func main() {
	goom()
	envDx()
	os.Exit(0)
}

//-----------------------------------------------------------------------------
