//-----------------------------------------------------------------------------
/*

Modules

*/
//-----------------------------------------------------------------------------

package core

import "fmt"

//-----------------------------------------------------------------------------
// Module Ports

type PortType int

const (
	PortType_Null        PortType = iota
	PortType_AudioBuffer          // audio buffers
	PortType_EventFloat           // event with float32 values
	PortType_EventMIDI            // event with MIDI data
)

type PortInfo struct {
	Name        string   // standard port name
	Description string   // description of port
	Ptype       PortType // port type
	Id          uint     // port ID: used as the ID for events on this port
}

//-----------------------------------------------------------------------------
// Module Information

type ModuleInfo struct {
	Name string     // module name
	In   []PortInfo // input ports
	Out  []PortInfo // input ports
}

// GetPortByName returns the module port information by port name.
func (mi *ModuleInfo) GetPortByName(name string) *PortInfo {
	// input ports
	for i := range mi.In {
		if name == mi.In[i].Name {
			return &mi.In[i]
		}
	}
	// output ports
	for i := range mi.Out {
		if name == mi.Out[i].Name {
			return &mi.Out[i]
		}
	}
	panic(fmt.Sprintf("no port named \"%s\" in module \"%s\"", name, mi.Name))
	return nil
}

// GetPortID returns the module port ID by port name.
func (mi *ModuleInfo) GetPortID(name string) uint {
	return mi.GetPortByName(name).Id
}

//-----------------------------------------------------------------------------
// Modules

type Module interface {
	Process(buf ...*Buf) // run the module dsp
	Event(e *Event)      // process an event
	Active() bool        // return true if the module has non-zero output
	Stop()               // stop the module
	Info() *ModuleInfo   // return the module information
}

//-----------------------------------------------------------------------------
