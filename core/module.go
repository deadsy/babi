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

// PortFuncType is a function used to send an event to the port of a module.
type PortFuncType func(m Module, e *Event)

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

// PortInfo contains the information describing a port.
type PortInfo struct {
	Name        string       // standard port name
	Description string       // description of port
	Ptype       PortType     // port type
	PortFunc    PortFuncType // port event function
}

// PortSet is a collection of ports.
type PortSet []PortInfo

//-----------------------------------------------------------------------------
// Module Information

// ModuleInfo contains the information describing a module.
type ModuleInfo struct {
	Name string                  // module name
	In   PortSet                 // input ports
	Out  PortSet                 // output ports
	n2p  map[string]PortFuncType // port name to port function mapping
}

// GetPortFunc returns the port function associated with the the port name.
func (mi *ModuleInfo) GetPortFunc(name string) PortFuncType {

	// build the name to port function map
	if mi.n2p == nil {
		// TODO detect duplicate port names
		mi.n2p = make(map[string]PortFuncType)
		// input ports
		for i := range mi.In {
			pf := mi.In[i].PortFunc
			if pf != nil {
				mi.n2p[mi.In[i].Name] = pf
			}
		}
		// output ports
		for i := range mi.Out {
			pf := mi.Out[i].PortFunc
			if pf != nil {
				mi.n2p[mi.Out[i].Name] = pf
			}
		}
	}
	// lookup the name
	if pf, ok := mi.n2p[name]; ok {
		return pf
	}
	log.Info.Printf("no port named \"%s\" in module \"%s\"", name, mi.Name)
	return nil
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

//-----------------------------------------------------------------------------
// Modules

// Module is the interface for an audio/event processing module.
type Module interface {
	Process(buf ...*Buf) // run the module dsp
	Active() bool        // return true if the module has non-zero output
	Stop()               // stop the module
	Info() *ModuleInfo   // return the module information
	Child() []Module     // return the child modules
}

// ModuleString returns a string for a tree of modules.
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

// ModuleName returns the name of a module
func ModuleName(m Module) string {
	return m.Info().Name
}

// ModuleStop calls Stop() for each module in a tree of modules.
func ModuleStop(m Module) {
	if m == nil {
		return
	}
	for _, c := range m.Child() {
		ModuleStop(c)
	}
	m.Stop()
}

//-----------------------------------------------------------------------------
