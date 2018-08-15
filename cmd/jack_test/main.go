//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package main

import (
	"github.com/deadsy/babi/jack"
	"github.com/deadsy/babi/log"
)

//-----------------------------------------------------------------------------

/*
 *
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

func main() {

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

var pcount int

func process(nframes uint32) int {
	if pcount%256 == 0 {
		//log.Info.Printf("%d\n", nframes)
	}
	pcount += 1

	return 0
}

func main() {

	client, status := jack.ClientOpen("jack_test", jack.NoStartServer)
	if status != 0 {
		log.Error.Printf("status: %s", jack.StrError(status))
		return
	}
	defer client.Close()

	if code := client.SetProcessCallback(process); code != 0 {
		log.Error.Printf("failed to set process callback: %s", jack.StrError(code))
		return
	}

	shutdown := make(chan struct{})
	client.OnShutdown(func() {
		log.Info.Printf("shutting down")
		close(shutdown)
	})

	// create two ports
	output_port1 := client.PortRegister("output_1", jack.DEFAULT_AUDIO_TYPE, jack.PortIsOutput, 0)
	output_port2 := client.PortRegister("output_2", jack.DEFAULT_AUDIO_TYPE, jack.PortIsOutput, 0)

	if code := client.Activate(); code != 0 {
		log.Error.Printf("failed to activate client: %s", jack.StrError(code))
		return
	}

	ports := client.GetPorts("", "", jack.PortIsInput|jack.PortIsPhysical)

	client.Connect(output_port1.GetName(), ports[0])
	client.Connect(output_port2.GetName(), ports[1])

	log.Info.Printf(client.GetName())
	<-shutdown
}

//-----------------------------------------------------------------------------
