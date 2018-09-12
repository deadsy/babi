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

// Info returns the module information.
func (m *seqModule) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "seq",
		In:   nil,
		Out:  nil,
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
		sm.pc -= 1
		return 1
	}
}

// OpNote returns a note operation.
func OpNote(channel, note, velocity, duration uint8) Op {
	return func(m *seqModule, sm *seqStateMachine) int {
		if sm.ostate == opStateInit {
			sm.duration = uint(duration)
			sm.ostate = opStateWait
			log.Info.Printf("note on %d (%d)", note, m.ticks)
			//seq_note_on(s, args);
		}
		sm.duration -= 1
		if sm.duration == 0 {
			// done
			sm.ostate = opStateInit
			log.Info.Printf("note off (%d)", m.ticks)
			//seq_note_off(s, args);
			return 1
		}
		// waiting...
		return 0
	}
}

// OpRest returns a rest operation.
func OpRest(duration uint8) Op {
	return func(m *seqModule, sm *seqStateMachine) int {
		if sm.ostate == opStateInit {
			sm.duration = uint(duration)
			sm.ostate = opStateWait
		}
		sm.duration -= 1
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
	if sm.sstate == seqStateRun {
		sm.pc += sm.prog[sm.pc](m, sm)
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
	synth       *core.Synth        // top-level synth
	beatsPerMin float32            // beats per minute
	secsPerTick float32            // seconds per tick
	tickError   float32            // current tick error
	ticks       uint               // full ticks
	sm          []*seqStateMachine // state machines
}

// NewSeq returns a basic sequencer module.
func NewSeq(s *core.Synth) core.Module {
	log.Info.Printf("")
	return &seqModule{
		synth: s,
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

// Event processes a module event.
func (m *seqModule) Event(e *core.Event) {
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *seqModule) Process(buf ...*core.Buf) {
}

// Active returns true if the module has non-zero output.
func (m *seqModule) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
