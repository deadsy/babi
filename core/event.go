//-----------------------------------------------------------------------------
/*

Events

*/
//-----------------------------------------------------------------------------

package core

//-----------------------------------------------------------------------------
// events

type EventType uint

const (
	Event_Null EventType = iota
	Event_MIDI
	Event_Ctrl
)

type Event struct {
	etype EventType   // event type
	info  interface{} // event information
}

func (e *Event) GetType() EventType {
	return e.etype
}

func (e *Event) GetMIDIEvent() *MIDIEvent {
	return e.info.(*MIDIEvent)
}

func (e *Event) GetCtrlEvent() *CtrlEvent {
	return e.info.(*CtrlEvent)
}

//-----------------------------------------------------------------------------
// MIDI events

type MIDIEventType uint

const (
	MIDIEvent_Null MIDIEventType = iota
	MIDIEvent_NoteOn
	MIDIEvent_NoteOff
	MIDIEvent_ControlChange
	MIDIEvent_PitchWheel
	MIDIEvent_PolyphonicAftertouch
	MIDIEvent_ProgramChange
	MIDIEvent_ChannelAftertouch
)

type MIDIEvent struct {
	etype  MIDIEventType
	status uint8 // message status byte
	arg0   uint8 // message byte 0
	arg1   uint8 // message byte 1
}

func (e *MIDIEvent) GetType() MIDIEventType {
	return e.etype
}

func (e *MIDIEvent) GetChannel() uint {
	return uint(e.status & 0xf)
}

func (e *MIDIEvent) GetNote() uint {
	return uint(e.arg0)
}

func (e *MIDIEvent) GetVelocity() uint {
	return uint(e.arg1)
}

//-----------------------------------------------------------------------------
// Control Event

type CtrlEventType uint

const (
	CtrlEvent_Null      CtrlEventType = iota
	CtrlEvent_NoteOn                  // trigger a note (key pressed)
	CtrlEvent_NoteOff                 // release a note (key released)
	CtrlEvent_Frequency               // set the frequency of an oscillator
)

type CtrlEvent struct {
	etype CtrlEventType
	val   float32
}

func (e *CtrlEvent) GetType() CtrlEventType {
	return e.etype
}

func (e *CtrlEvent) GetVal() float32 {
	return e.val
}

//-----------------------------------------------------------------------------
