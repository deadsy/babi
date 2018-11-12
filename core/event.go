//-----------------------------------------------------------------------------
/*

Events

*/
//-----------------------------------------------------------------------------

package core

import (
	"fmt"
)

//-----------------------------------------------------------------------------
// generic events

// Event is a generic event.
type Event struct {
	info interface{} // event information
}

// NewEvent returns an event.
func NewEvent(info interface{}) *Event {
	return &Event{info}
}

// String returns a descriptive string for the event.
func (e *Event) String() string {
	me := e.GetEventMIDI()
	if me != nil {
		return me.String()
	}
	fe := e.GetEventFloat()
	if fe != nil {
		return fe.String()
	}
	ie := e.GetEventInt()
	if ie != nil {
		return ie.String()
	}
	return "unknown event"
}

// SendEvent sends an event to a named port on a module.
func SendEvent(m Module, name string, e *Event) {
	portFunc := m.Info().GetPortFunc(name)
	if portFunc != nil {
		portFunc(m, e)
	}
}

//-----------------------------------------------------------------------------
// MIDI Events

// EventTypeMIDI is the MIDI event type.
// See: https://www.midi.org/specifications
type EventTypeMIDI uint

// EventTypeMIDI enumeration.
const (
	EventMIDINull                 EventTypeMIDI = 0
	EventMIDINoteOn                             = midiStatusNoteOn
	EventMIDINoteOff                            = midiStatusNoteOff
	EventMIDIControlChange                      = midiStatusControlChange
	EventMIDIPitchWheel                         = midiStatusPitchWheel
	EventMIDIPolyphonicAftertouch               = midiStatusPolyphonicAftertouch
	EventMIDIProgramChange                      = midiStatusProgramChange
	EventMIDIChannelAftertouch                  = midiStatusChannelAftertouch
)

var midiEventType2String = map[EventTypeMIDI]string{
	EventMIDINull:                 "null",
	EventMIDINoteOn:               "note_on",
	EventMIDINoteOff:              "note_off",
	EventMIDIControlChange:        "control_change",
	EventMIDIPitchWheel:           "pitch_wheel",
	EventMIDIPolyphonicAftertouch: "polyphonic_aftertouch",
	EventMIDIProgramChange:        "program_change",
	EventMIDIChannelAftertouch:    "channel_aftertouch",
}

// EventMIDI is an event with MIDI data.
type EventMIDI struct {
	etype  EventTypeMIDI
	status uint8 // message status byte
	arg0   uint8 // message byte 0
	arg1   uint8 // message byte 1
}

// NewEventMIDI returns a new MIDI event.
func NewEventMIDI(etype EventTypeMIDI, status, arg0, arg1 uint8) *Event {
	return NewEvent(&EventMIDI{etype, status, arg0, arg1})
}

// String returns a descriptive string for the MIDI event.
func (e *EventMIDI) String() string {
	descr := midiEventType2String[e.etype]
	switch e.GetType() {
	case EventMIDINoteOn, EventMIDINoteOff:
		return fmt.Sprintf("%s ch %d note %d vel %d", descr, e.GetChannel(), e.GetNote(), e.GetVelocity())
	case EventMIDIControlChange:
		return fmt.Sprintf("%s ch %d ctrl %d val %d", descr, e.GetChannel(), e.GetCtrlNum(), e.GetCtrlVal())
	case EventMIDIPitchWheel:
		return fmt.Sprintf("%s ch %d val %d", descr, e.GetChannel(), e.GetPitchWheel())
		//case EventMIDIPolyphonicAftertouch:
		//case EventMIDIProgramChange:
		//case EventMIDIChannelAftertouch:
	}
	return fmt.Sprintf("%s status %02x arg0 %02x arg1 %02x", descr, e.status, e.arg0, e.arg1)
}

// GetEventMIDI returns a MIDI event (or nil).
func (e *Event) GetEventMIDI() *EventMIDI {
	if me, ok := e.info.(*EventMIDI); ok {
		return me
	}
	return nil
}

// GetEventMIDIChannel returns the MIDI event for the MIDI channel.
func (e *Event) GetEventMIDIChannel(ch uint8) *EventMIDI {
	me := e.GetEventMIDI()
	if me != nil && me.GetChannel() == ch {
		return me
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
	return uint16(e.arg1)<<7 | uint16(e.arg0)
}

//-----------------------------------------------------------------------------
// Float Events

// EventFloat is an event with a 32-bit floating point value.
type EventFloat struct {
	Val float32
}

// NewEventFloat returns a new control event.
func NewEventFloat(val float32) *Event {
	return NewEvent(&EventFloat{val})
}

// String returns a descriptive string for the float event.
func (e *EventFloat) String() string {
	return fmt.Sprintf("val %f", e.Val)
}

// GetEventFloat returns the float event.
func (e *Event) GetEventFloat() *EventFloat {
	if fe, ok := e.info.(*EventFloat); ok {
		return fe
	}
	return nil
}

// SendEventFloat sends a float event to a named port on a module.
func SendEventFloat(m Module, name string, val float32) {
	portFunc := m.Info().GetPortFunc(name)
	if portFunc != nil {
		portFunc(m, NewEventFloat(val))
	}
}

//-----------------------------------------------------------------------------
// Integer Events

// EventInt is an event with a 32-bit signed integer value.
type EventInt struct {
	Val int
}

// NewEventInt returns a new integer event.
func NewEventInt(val int) *Event {
	return NewEvent(&EventInt{val})
}

// String returns a descriptive string for the integer event.
func (e *EventInt) String() string {
	return fmt.Sprintf("val %d", e.Val)
}

// GetEventInt returns the integer event.
func (e *Event) GetEventInt() *EventInt {
	if ie, ok := e.info.(*EventInt); ok {
		return ie
	}
	return nil
}

// SendEventInt sends a integer event to a named port on a module.
func SendEventInt(m Module, name string, val int) {
	portFunc := m.Info().GetPortFunc(name)
	if portFunc != nil {
		portFunc(m, NewEventInt(val))
	}
}

//-----------------------------------------------------------------------------
