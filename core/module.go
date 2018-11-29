//-----------------------------------------------------------------------------
/*

Modules

*/
//-----------------------------------------------------------------------------

package core

import (
	"fmt"
	"strings"
	"sync"
)

//-----------------------------------------------------------------------------
// Module Ports

// PortFuncType is a function used to send an event to the input port of a module.
type PortFuncType func(m Module, e *Event)

// PortType represents the type of data sent or received on a module port.
type PortType int

// PortType enumeration.
const (
	PortTypeNull  PortType = iota
	PortTypeAudio          // audio buffers
	PortTypeFloat          // event with float32 values
	PortTypeInt            // event with integer values
	PortTypeBool           // event with boolean values
	PortTypeMIDI           // event with MIDI data
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

// dstPort stores the module and port function for sending an event
type dstPort struct {
	module   Module       // destination module
	portFunc PortFuncType // destination port function
}

// ModuleInfo describes a modules ports and connectivity.
type ModuleInfo struct {
	Name   string                  // module name
	In     PortSet                 // input ports
	Out    PortSet                 // output ports
	synth  *Synth                  // top-level synth
	inMap  map[string]PortFuncType // map input port names to port functions
	outMap map[string][]dstPort    // map output port names to input ports of other modules
}

// getPortFunc returns the port function for the input port name.
func (mi *ModuleInfo) getPortFunc(name string) PortFuncType {
	if pf, ok := mi.inMap[name]; ok {
		return pf
	}
	//log.Info.Printf("module \"%s\" has no port named \"%s\"", mi.Name, name)
	return nil
}

//-----------------------------------------------------------------------------

// numPortsByType return the number of ports within a set matching a type.
func (ps PortSet) numPortsByType(ptype PortType) int {
	var n int
	for _, pi := range ps {
		if pi.Ptype == ptype {
			n++
		}
	}
	return n
}

// numPortsByName returns the number of ports within a set matching a name.
func (ps PortSet) numPortsByName(name string) int {
	var n int
	for _, pi := range ps {
		if pi.Name == name {
			n++
		}
	}
	return n
}

//-----------------------------------------------------------------------------

// Connect source/destination module event ports.
func Connect(s Module, sname string, d Module, dname string) {

	si := s.Info()
	di := d.Info()

	// check output on source module
	n := si.Out.numPortsByName(sname)
	if n != 1 {
		panic(fmt.Sprintf("module \"%s\" must have one output port with name \"%s\"", si.Name, sname))
	}
	// check input on destination module
	n = di.In.numPortsByName(dname)
	if n != 1 {
		panic(fmt.Sprintf("module \"%s\" must have one input port with name \"%s\"", di.Name, dname))
	}
	// TODO check port types

	// destination port function
	pf := di.getPortFunc(dname)
	if pf == nil {
		return
	}
	// add it to the output port mapping for this source
	si.outMap[sname] = append(si.outMap[sname], dstPort{d, pf})
}

// EventOut sends an event from the named output port of a module.
// The event will be sent to input ports connected to the output port.
func EventOut(m Module, name string, e *Event) {
	mi := m.Info()
	if dstPorts, ok := mi.outMap[name]; ok {
		for i := range dstPorts {
			dstPorts[i].portFunc(dstPorts[i].module, e)
		}
	}
}

//-----------------------------------------------------------------------------
// Modules

// Module is the interface for an audio/event processing module.
type Module interface {
	Process(buf ...*Buf) // run the module dsp
	Active() bool        // return true if the module should be run for dsp output
	Stop()               // stop the module
	Info() *ModuleInfo   // return module information
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

var globalID uint
var globalIDLock sync.Mutex

func getID() uint {
	globalIDLock.Lock()
	if globalID == 0 {
		globalID = 1
	}
	id := globalID
	globalID++
	globalIDLock.Unlock()
	return id
}

// GenerateID returns a system unique ID string (with name prefix).
func GenerateID(name string) string {
	return fmt.Sprintf("%s%d", name, getID())
}

//-----------------------------------------------------------------------------
