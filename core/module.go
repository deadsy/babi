//-----------------------------------------------------------------------------
/*

Modules

*/
//-----------------------------------------------------------------------------

package core

//-----------------------------------------------------------------------------

type Module interface {
	Process(buf []*Buf) // run the module dsp
	Event(e *Event)     // process an event
	Active() bool       // return true if the module has non-zero output
	Stop()              // stop the module
	Ports() []PortInfo  // return the module port information
}

//-----------------------------------------------------------------------------
// module ports

type PortType int

const (
	PortType_Null PortType = iota
	PortType_Buf           // audio buffers
	PortType_Ctrl          // control events
	PortType_MIDI          // midi events
)

type PortDirn int

const (
	PortDirn_Null PortDirn = iota
	PortDirn_In            // input
	PortDirn_Out           // output
)

type PortInfo struct {
	Label       string      // small label for port
	Description string      // description of port
	Ptype       PortType    // port type
	Dirn        PortDirn    // port direction
	Info        interface{} // extra port info
}

//-----------------------------------------------------------------------------
