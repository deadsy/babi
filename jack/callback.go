//-----------------------------------------------------------------------------
/*

Golang wrapper for jackd2.

*/
//-----------------------------------------------------------------------------

package jack

import "C"
import "unsafe"

//-----------------------------------------------------------------------------

type ProcessCallback func(uint32) int
type BufferSizeCallback func(uint32) int
type SampleRateCallback func(uint32) int
type PortRegistrationCallback func(PortId, bool)
type PortRenameCallback func(PortId, string, string)
type PortConnectCallback func(PortId, PortId, bool)
type ClientRegistrationCallback func()
type FreewheelCallback func()
type GraphOrderCallback func() int
type XrunCallback func() int
type ShutdownCallback func()
type InfoShutdownCallback func()
type ErrorFunction func(string)
type InfoFunction func(string)

//export goProcess
func goProcess(nframes uint, arg unsafe.Pointer) int {
	client := (*C.struct__jack_client)(arg)
	return clientMap[client].processCallback(uint32(nframes))
}

//export goBufferSize
func goBufferSize(nframes uint, arg unsafe.Pointer) int {
	client := (*C.struct__jack_client)(arg)
	return clientMap[client].bufferSizeCallback(uint32(nframes))
}

//export goSampleRate
func goSampleRate(nframes uint, arg unsafe.Pointer) int {
	client := (*C.struct__jack_client)(arg)
	return clientMap[client].sampleRateCallback(uint32(nframes))
}

//export goPortRegistration
func goPortRegistration(port uint, register int, arg unsafe.Pointer) {
	client := (*C.struct__jack_client)(arg)
	clientMap[client].portRegistrationCallback(PortId(port), register != 0)
}

//export goPortRename
func goPortRename(port uint, oldName, newName *C.char, arg unsafe.Pointer) {
	client := (*C.struct__jack_client)(arg)
	clientMap[client].portRenameCallback(PortId(port), C.GoString(oldName), C.GoString(newName))
}

//export goPortConnect
func goPortConnect(aport, bport uint, connect int, arg unsafe.Pointer) {
	client := (*C.struct__jack_client)(arg)
	clientMap[client].portConnectCallback(PortId(aport), PortId(bport), connect != 0)
}

//export goClientRegistration
func goClientRegistration(name *C.char, reg int, arg unsafe.Pointer) {
	client := (*C.struct__jack_client)(arg)
	clientMap[client].clientRegistrationCallback()
}

//export goFreewheel
func goFreewheel(starting int, arg unsafe.Pointer) {
	client := (*C.struct__jack_client)(arg)
	clientMap[client].freewheelCallback()
}

//export goGraphOrder
func goGraphOrder(arg unsafe.Pointer) int {
	client := (*C.struct__jack_client)(arg)
	return clientMap[client].graphOrderCallback()
}

//export goXrun
func goXrun(arg unsafe.Pointer) int {
	client := (*C.struct__jack_client)(arg)
	return clientMap[client].xrunCallback()
}

//export goShutdown
func goShutdown(arg unsafe.Pointer) {
	client := (*C.struct__jack_client)(arg)
	clientMap[client].shutdownCallback()
}

//export goInfoShutdown
func goInfoShutdown(code uint, reason *C.char, arg unsafe.Pointer) {
	client := (*C.struct__jack_client)(arg)
	clientMap[client].infoShutdownCallback()
}

//export goErrorFunction
func goErrorFunction(msg *C.char) {
	if errorFunction != nil {
		errorFunction(C.GoString(msg))
	}
}

//export goInfoFunction
func goInfoFunction(msg *C.char) {
	if infoFunction != nil {
		infoFunction(C.GoString(msg))
	}
}

//-----------------------------------------------------------------------------
