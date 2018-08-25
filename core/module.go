//-----------------------------------------------------------------------------
/*

Modules

*/
//-----------------------------------------------------------------------------

package core

//-----------------------------------------------------------------------------

type Module interface {
	Process(in, out []*Buf) // run the module dsp
	Event(e *Event)         // process an event
	Active() bool           // return true if the module has non-zero output
	Stop()                  // stop the module
}

//-----------------------------------------------------------------------------
