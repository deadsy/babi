//-----------------------------------------------------------------------------
/*

Events

*/
//-----------------------------------------------------------------------------

package core

import "fmt"

//-----------------------------------------------------------------------------
// general event

type eventType uint

const (
	eventTypeNull eventType = iota
	eventTypeMIDI
	eventTypeFloat
	eventTypeInt
)

var eventType2String = map[eventType]string{
	eventTypeNull:  "null",
	eventTypeMIDI:  "midi",
	eventTypeFloat: "float",
	eventTypeInt:   "int",
}

type Event struct {
	etype eventType   // event type
	info  interface{} // event information
}

// NewEvent returns a new event.
func NewEvent(etype eventType, info interface{}) *Event {
	return &Event{etype, info}
}

// String returns a descriptive string for the event.
func (e *Event) String() string {
	var s string
	switch e.etype {
	case eventTypeNull:
		s = "null"
	case eventTypeMIDI:
		s = fmt.Sprintf("%s", e.info.(*EventMIDI))
	case eventTypeFloat:
		s = fmt.Sprintf("%s", e.info.(*EventFloat))
	case eventTypeInt:
		s = fmt.Sprintf("%s", e.info.(*EventInt))
	default:
		s = "unknown"
	}
	return fmt.Sprintf("%s: %s", eventType2String[e.etype], s)
}

//-----------------------------------------------------------------------------
// MIDI Events

type EventTypeMIDI uint

const (
	EventMIDI_Null EventTypeMIDI = iota
	EventMIDI_NoteOn
	EventMIDI_NoteOff
	EventMIDI_ControlChange
	EventMIDI_PitchWheel
	EventMIDI_PolyphonicAftertouch
	EventMIDI_ProgramChange
	EventMIDI_ChannelAftertouch
)

var midiEventType2String = map[EventTypeMIDI]string{
	EventMIDI_Null:                 "null",
	EventMIDI_NoteOn:               "note_on",
	EventMIDI_NoteOff:              "note_off",
	EventMIDI_ControlChange:        "control_change",
	EventMIDI_PitchWheel:           "pitch_wheel",
	EventMIDI_PolyphonicAftertouch: "polyphonic_aftertouch",
	EventMIDI_ProgramChange:        "program_change",
	EventMIDI_ChannelAftertouch:    "channel_aftertouch",
}

type EventMIDI struct {
	etype  EventTypeMIDI
	status uint8 // message status byte
	arg0   uint8 // message byte 0
	arg1   uint8 // message byte 1
}

// NewEventMIDI returns a new MIDI event.
func NewEventMIDI(etype EventTypeMIDI, status, arg0, arg1 uint8) *Event {
	return NewEvent(eventTypeMIDI, &EventMIDI{etype, status, arg0, arg1})
}

// String returns a descriptive string for the MIDI event.
func (e *EventMIDI) String() string {
	descr := midiEventType2String[e.etype]
	switch e.GetType() {
	case EventMIDI_NoteOn, EventMIDI_NoteOff:
		return fmt.Sprintf("%s ch %d note %d vel %d", descr, e.GetChannel(), e.GetNote(), e.GetVelocity())
	case EventMIDI_ControlChange:
		return fmt.Sprintf("%s ch %d ctrl %d val %d", descr, e.GetChannel(), e.GetCtrlNum(), e.GetCtrlVal())
	case EventMIDI_PitchWheel:
		return fmt.Sprintf("%s ch %d val %d", descr, e.GetChannel(), e.GetPitchWheel())
		//case EventMIDI_PolyphonicAftertouch:
		//case EventMIDI_ProgramChange:
		//case EventMIDI_ChannelAftertouch:
	}
	return fmt.Sprintf("%s status %02x arg0 %02x arg1 %02x", midiEventType2String[e.etype], e.status, e.arg0, e.arg1)
}

// GetEventMIDIChannel returns the MIDI event for the MIDI channel.
func (e *Event) GetEventMIDIChannel(ch uint8) *EventMIDI {
	if me, ok := e.info.(*EventMIDI); ok {
		if me.GetChannel() == ch {
			return me
		}
	}
	return nil
}

// GetType returns the MIDI event type.
func (e *EventMIDI) GetType() EventTypeMIDI {
	return e.etype
}

// GetChannel returns the MIDI channel number.
func (e *EventMIDI) GetChannel() uint8 {
	return e.status & 0xf
}

// GetNote returns the MIDI note value.
func (e *EventMIDI) GetNote() uint8 {
	return e.arg0
}

// GetCtrlNum returns the MIDI control number.
func (e *EventMIDI) GetCtrlNum() uint8 {
	return e.arg0
}

// GetCtrlVal returns the MIDI control value.
func (e *EventMIDI) GetCtrlVal() uint8 {
	return e.arg1
}

// GetVelocity returns the MIDI note velocity.
func (e *EventMIDI) GetVelocity() uint8 {
	return e.arg1
}

// GetPitchWheel returns the MIDI pitch wheel value.
func (e *EventMIDI) GetPitchWheel() uint16 {
	return uint16(e.arg1<<7) | uint16(e.arg0)
}

//-----------------------------------------------------------------------------
// Float Events

type EventFloat struct {
	Id  PortId
	Val float32
}

// NewEventFloat returns a new control event.
func NewEventFloat(id PortId, val float32) *Event {
	return NewEvent(eventTypeFloat, &EventFloat{id, val})
}

// String returns a descriptive string for the float event.
func (e *EventFloat) String() string {
	return fmt.Sprintf("id %d val %f", e.Id, e.Val)
}

// GetEventFloat returns the float event.
func (e *Event) GetEventFloat() *EventFloat {
	if fe, ok := e.info.(*EventFloat); ok {
		return fe
	}
	return nil
}

// SendEventFloatName sends a float event to a named port on a module.
func SendEventFloatName(m Module, name string, val float32) {
	m.Event(NewEventFloat(m.Info().GetPortByName(name).Id, val))
}

// SendEventFloatID sends a float event to a port ID on a module.
func SendEventFloatID(m Module, id PortId, val float32) {
	m.Event(NewEventFloat(id, val))
}

//-----------------------------------------------------------------------------
// Integer Events

type EventInt struct {
	Id  PortId
	Val int
}

// NewEventInt returns a new integer event.
func NewEventInt(id PortId, val int) *Event {
	return NewEvent(eventTypeInt, &EventInt{id, val})
}

// String returns a descriptive string for the integer event.
func (e *EventInt) String() string {
	return fmt.Sprintf("id %d val %d", e.Id, e.Val)
}

// GetEventInt returns the integer event.
func (e *Event) GetEventInt() *EventInt {
	if ie, ok := e.info.(*EventInt); ok {
		return ie
	}
	return nil
}

// SendEventIntName sends a integer event to a named port on a module.
func SendEventIntName(m Module, name string, val int) {
	m.Event(NewEventInt(m.Info().GetPortByName(name).Id, val))
}

// SendEventIntID sends a integer event to a port ID on a module.
func SendEventIntID(m Module, id PortId, val int) {
	m.Event(NewEventInt(id, val))
}

//-----------------------------------------------------------------------------
