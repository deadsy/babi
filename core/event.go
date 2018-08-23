//-----------------------------------------------------------------------------
/*

Events

*/
//-----------------------------------------------------------------------------

package core

//-----------------------------------------------------------------------------
// general event

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

// NewEvent returns a new event.
func NewEvent(etype EventType, info interface{}) *Event {
	return &Event{etype, info}
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

// NewMIDIEvent returns a new MIDI event.
func NewMIDIEvent(etype MIDIEventType, status, arg0, arg1 uint8) *Event {
	return NewEvent(Event_MIDI, &MIDIEvent{etype, status, arg0, arg1})
}

// GetType returns the MIDI event type.
func (e *MIDIEvent) GetType() MIDIEventType {
	return e.etype
}

// GetChannel returns the MIDI channel number.
func (e *MIDIEvent) GetChannel() uint8 {
	return e.status & 0xf
}

// GetNote returns the MIDI note value.
func (e *MIDIEvent) GetNote() uint8 {
	return e.arg0
}

// GetVelocity returns the MIDI note velocity.
func (e *MIDIEvent) GetVelocity() uint8 {
	return e.arg1
}

// GetPitchWheel returns the MIDI pitch wheel value.
func (e *MIDIEvent) GetPitchWheel() uint16 {
	return uint16(e.arg1<<7) | uint16(e.arg0)
}

//-----------------------------------------------------------------------------
// Control Event

type CtrlEventType uint

const (
	CtrlEvent_Null      CtrlEventType = iota
	CtrlEvent_NoteOn                  // trigger a note (key pressed)
	CtrlEvent_NoteOff                 // release a note (key released)
	CtrlEvent_Frequency               // set an oscillator frequency
	CtrlEvent_Attenuate               // set an attenuation level
)

type CtrlEvent struct {
	etype CtrlEventType
	val   float32
}

// NewCtrlEvent returns a new control event.
func NewCtrlEvent(etype CtrlEventType, val float32) *Event {
	return NewEvent(Event_Ctrl, &CtrlEvent{etype, val})
}

// GetType returns the type of a control event.
func (e *CtrlEvent) GetType() CtrlEventType {
	return e.etype
}

// GetVal returns the value of a control event.
func (e *CtrlEvent) GetVal() float32 {
	return e.val
}

//-----------------------------------------------------------------------------
