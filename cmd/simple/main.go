//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package main

import (
	"fmt"

	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/patches"
)

//-----------------------------------------------------------------------------

/*

var channels int = 2

var PortsIn []*jack.Port
var PortsOut []*jack.Port

func process(nframes uint32) int {
	for i, in := range PortsIn {
		samplesIn := in.GetBuffer(nframes)
		samplesOut := PortsOut[i].GetBuffer(nframes)
		for i2, sample := range samplesIn {
			samplesOut[i2] = sample
		}
	}
	return 0
}

func main_x() {

	client, status := jack.ClientOpen("Go Passthrough", jack.NoStartServer)
	if status != 0 {
		fmt.Println("Status:", jack.StrError(status))
		return
	}
	defer client.Close()

	if code := client.SetProcessCallback(process); code != 0 {
		fmt.Println("Failed to set process callback:", jack.StrError(code))
		return
	}
	shutdown := make(chan struct{})
	client.OnShutdown(func() {
		fmt.Println("Shutting down")
		close(shutdown)
	})

	if code := client.Activate(); code != 0 {
		fmt.Println("Failed to activate client:", jack.StrError(code))
		return
	}

	for i := 0; i < channels; i++ {
		portIn := client.PortRegister(fmt.Sprintf("in_%d", i), jack.DEFAULT_AUDIO_TYPE, jack.PortIsInput, 0)
		PortsIn = append(PortsIn, portIn)
	}
	for i := 0; i < channels; i++ {
		portOut := client.PortRegister(fmt.Sprintf("out_%d", i), jack.DEFAULT_AUDIO_TYPE, jack.PortIsOutput, 0)
		PortsOut = append(PortsOut, portOut)
	}

	fmt.Println(client.GetName())
	<-shutdown
}

*/

//-----------------------------------------------------------------------------

/*

func simple(pa *pulse.PulseMainLoop) {

	ctx := pa.NewContext("default", 0)
	if ctx == nil {
		fmt.Println("Failed to create a new context")
		return
	}
	defer ctx.Dispose()
	st := ctx.NewStream("default", &pulse.PulseSampleSpec{Format: pulse.SAMPLE_FLOAT32LE, Rate: core.AUDIO_FS, Channels: 1})
	if st == nil {
		fmt.Println("Failed to create a new stream")
		return
	}
	defer st.Dispose()
	st.ConnectToSink()

	p := patches.NewSimple()

	for {
		core.AudioClear()
		p.Process()
		st.Write(core.AudioGetL(), pulse.SEEK_RELATIVE)
	}

}

*/

func main() {

	audio, err := core.NewPulse()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	defer audio.Close()

	b := core.NewBabi(audio)
	b.AddPatch(patches.NewSimple(b))
	b.Run()
}

//-----------------------------------------------------------------------------
