//-----------------------------------------------------------------------------
/*

Modules

*/
//-----------------------------------------------------------------------------

package core

//-----------------------------------------------------------------------------

type Module interface {
	Process(buf ...*Buf) // run the module dsp
	Event(e *Event)      // process an event
	Active() bool        // return true if the module has non-zero output
	Stop()               // stop the module
	Info() *ModuleInfo   // return the module information
}

type ModuleInfo struct {
	In  []PortInfo // input ports
	Out []PortInfo // input ports
}

//-----------------------------------------------------------------------------
// module ports

type PortType int

const (
	PortType_Null         PortType = iota
	PortType_AudioBuffer           // audio buffers
	PortType_EventFloat32          // event with float32 values
	PortType_EventMIDI             // event with MIDI data
)

type PortInfo struct {
	Label       string   // short label for port
	Description string   // description of port
	Ptype       PortType // port type
}

//-----------------------------------------------------------------------------
