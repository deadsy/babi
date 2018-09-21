//-----------------------------------------------------------------------------
/*

Modules

*/
//-----------------------------------------------------------------------------

package core

import (
	"fmt"
	"strings"

	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------
// Module Ports

// PortType represents the type of data sent or received on a module port.
type PortType int

// PortType enumeration.
const (
	PortTypeNull        PortType = iota
	PortTypeAudioBuffer          // audio buffers
	PortTypeFloat                // event with float32 values
	PortTypeInt                  // event with integer values
	PortTypeMIDI                 // event with MIDI data
)

// PortID is a numeric identifier for a module port.
type PortID uint

// PortInfo contains the information describing a port.
type PortInfo struct {
	Name        string   // standard port name
	Description string   // description of port
	Ptype       PortType // port type
	ID          PortID   // numeric port id
}

// PortSet is a collection of ports.
type PortSet []PortInfo

//-----------------------------------------------------------------------------
// Module Information

// ModuleInfo contains the information describing a module.
type ModuleInfo struct {
	Name string  // module name
	In   PortSet // input ports
	Out  PortSet // output ports
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
	log.Info.Printf("no port named \"%s\" in module \"%s\"", name, mi.Name)
	return nil
}

// GetPortID returns the module port ID by port name.
func (mi *ModuleInfo) GetPortID(name string) PortID {
	pi := mi.GetPortByName(name)
	if pi != nil {
		return pi.ID
	}
	return 0
}

// numPorts return the number of ports within a set matching a specific type.
func (ps PortSet) numPorts(ptype PortType) int {
	var count int
	for _, pi := range ps {
		if pi.Ptype == ptype {
			count++
		}
	}
	return count
}

// CheckIO checks a module for the required type/number of IO ports.
func (mi *ModuleInfo) CheckIO(midiIn, audioIn, audioOut int) error {
	var n int
	n = mi.In.numPorts(PortTypeMIDI)
	if n != midiIn {
		return fmt.Errorf("%s needs %d MIDI inputs (has %d)", mi.Name, midiIn, n)
	}
	n = mi.In.numPorts(PortTypeAudioBuffer)
	if n != audioIn {
		return fmt.Errorf("%s needs %d audio inputs (has %d)", mi.Name, audioIn, n)
	}
	n = mi.Out.numPorts(PortTypeAudioBuffer)
	if n != audioOut {
		return fmt.Errorf("%s needs %d audio outputs (has %d)", mi.Name, audioOut, n)
	}
	return nil
}

//-----------------------------------------------------------------------------
// Modules

// Module is the interface for an audio/event processing module.
type Module interface {
	Process(buf ...*Buf) // run the module dsp
	Event(e *Event)      // process an event
	Active() bool        // return true if the module has non-zero output
	Stop()               // stop the module
	Info() *ModuleInfo   // return the module information
	Child() []Module     // return the child modules
}

// ModuleString returns the string for the module tree of this module.
func ModuleString(m Module) string {
	mi := m.Info()
	children := m.Child()
	if len(children) != 0 {
		s := make([]string, len(children))
		for i, c := range children {
			s[i] = ModuleString(c)
		}
		return fmt.Sprintf("%s (%s)", mi.Name, strings.Join(s, " "))
	}
	return fmt.Sprintf("%s", mi.Name)
}

//-----------------------------------------------------------------------------
