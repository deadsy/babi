//-----------------------------------------------------------------------------
/*

Basic Sequencer

*/
//-----------------------------------------------------------------------------

package seq

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

const ticksPerBeat = 16

//-----------------------------------------------------------------------------

// Info returns the module information.
func (m *seqModule) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "seq",
		In: []core.PortInfo{
			{"bpm", "beats per minute", core.PortTypeFloat, seqPortBpm},
			{"ctrl", "control", core.PortTypeInt, seqPortCtrl},
		},
		Out: []core.PortInfo{
			{"midi_out", "midi output", core.PortTypeMIDI, nil},
		},
	}
}

//-----------------------------------------------------------------------------
// operations

// Op is a sequencer operation function.
type Op func(m *seqModule, sm *seqStateMachine) int

// OpNOP returns a nop operation.
func OpNOP() Op {
	return func(m *seqModule, sm *seqStateMachine) int {
		return 1
	}
}

// OpLoop returns a loop operation.
func OpLoop() Op {
	return func(m *seqModule, sm *seqStateMachine) int {
		log.Info.Printf("loop (%d)", m.ticks)
		sm.pc = -1
		return 1
	}
}

// OpNote returns a note operation.
func OpNote(channel, note, velocity uint8, duration uint) Op {
	return func(m *seqModule, sm *seqStateMachine) int {
		if sm.ostate == opStateInit {
			sm.duration = duration
			sm.ostate = opStateWait
			log.Info.Printf("note on %d (%d)", note, m.ticks)
			m.synth.PushEvent(core.NewEventMIDI(core.EventMIDINoteOn, channel, note, velocity))
		}
		sm.duration--
		if sm.duration == 0 {
			// done
			sm.ostate = opStateInit
			log.Info.Printf("note off (%d)", m.ticks)
			m.synth.PushEvent(core.NewEventMIDI(core.EventMIDINoteOff, channel, note, 0))
			return 1
		}
		// waiting...
		return 0
	}
}

// OpRest returns a rest operation.
func OpRest(duration uint) Op {
	return func(m *seqModule, sm *seqStateMachine) int {
		if sm.ostate == opStateInit {
			sm.duration = duration
			sm.ostate = opStateWait
		}
		sm.duration--
		if sm.duration == 0 {
			// done
			sm.ostate = opStateInit
			return 1
		}
		// waiting...
		return 0
	}
}

func (m *seqModule) tick(sm *seqStateMachine) {
	// auto stop zero length programs
	if len(sm.prog) == 0 {
		sm.sstate = seqStateStop
	}
	// run the program
	if sm.sstate == seqStateRun {
		n := sm.prog[sm.pc](m, sm)
		sm.pc += n
	}
}

//-----------------------------------------------------------------------------

type seqState int

const (
	seqStateStop seqState = iota // initial state
	seqStateRun
)

type opState int

const (
	opStateInit opState = iota // initial state
	opStateWait
)

// per sequence state machine
type seqStateMachine struct {
	prog     []Op     // program operations
	pc       int      // program counter
	sstate   seqState // sequencer state
	ostate   opState  // operation state
	duration uint     // operation duration
}

type seqModule struct {
	synth       *core.Synth      // top-level synth
	secsPerTick float32          // seconds per tick
	tickError   float32          // current tick error
	ticks       uint             // full ticks
	sm          *seqStateMachine // state machine
}

// NewSequencer returns a basic sequencer module.
func NewSequencer(s *core.Synth, prog []Op) core.Module {
	log.Info.Printf("")
	return &seqModule{
		synth: s,
		sm:    &seqStateMachine{prog: prog},
	}
}

// Child returns the child modules of this module.
func (m *seqModule) Child() []core.Module {
	return nil
}

// Stop performs any cleanup of a module.
func (m *seqModule) Stop() {
}

//-----------------------------------------------------------------------------
// Events

// Sequencer control values.
const (
	CtrlStop  = iota // stop the sequencer
	CtrlStart        // start the sequencer
	CtrlReset        // reset the sequencer
)

func seqPortBpm(cm core.Module, e *core.Event) {
	m := cm.(*seqModule)
	bpm := core.Clamp(e.GetEventFloat().Val, core.MinBeatsPerMin, core.MaxBeatsPerMin)
	log.Info.Printf("set bpm %f", bpm)
	m.secsPerTick = core.SecsPerMin / (bpm * ticksPerBeat)
}

func seqPortCtrl(cm core.Module, e *core.Event) {
	m := cm.(*seqModule)
	ctrl := e.GetEventInt().Val
	switch ctrl {
	case CtrlStop: // stop the sequencer
		log.Info.Printf("ctrl stop")
		m.sm.sstate = seqStateStop
	case CtrlStart: // start the sequencer
		log.Info.Printf("ctrl start")
		m.sm.sstate = seqStateRun
	case CtrlReset: // reset the sequencer
		log.Info.Printf("ctrl reset")
		m.sm.sstate = seqStateStop
		m.sm.ostate = opStateInit
		m.sm.pc = 0
	default:
		log.Info.Printf("unknown control value %d", ctrl)
	}
}

// Event processes a module event.
func (m *seqModule) Event(e *core.Event) {
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *seqModule) Process(buf ...*core.Buf) {
	// This routine is being used as a periodic call for timed event generation.
	// The sequencer does not process audio buffers.

	// The desired BPM will generally not correspond to an integral number
	// of audio blocks, so accumulate an error and tick when needed.
	// ie- Bresenham style.
	m.tickError += core.SecsPerAudioBuffer
	if m.tickError > m.secsPerTick {
		m.tickError -= m.secsPerTick
		m.ticks++
		// tick the state machine
		m.tick(m.sm)
	}
}

// Active returns true if the module has non-zero output.
func (m *seqModule) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
