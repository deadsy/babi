//-----------------------------------------------------------------------------
/*

Basic Sequencer

*/
//-----------------------------------------------------------------------------

package seq

import (
	"fmt"

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
// opcodes

const (
	opNOP  = iota // no operation
	opLoop        // return to beginning
	opNote        // note on/off
	opRest        // rest
)

// no operation, (op)
func (m *seqModule) opNOP(sm *seqStateMachine) int {
	return 1
}

// return to beginning, (op)
func (m *seqModule) opLoop(sm *seqStateMachine) int {
	sm.pc = -1
	return 1
}

// note on/off, (op, channel, note, velocity, duration)
func (m *seqModule) opNote(sm *seqStateMachine) int {
	/*
		struct note_args *args = (struct note_args *)&m->prog[m->pc];
		if (m->op_state == O_STATE_INIT) {
			// init
			m->duration = args->dur;
			m->op_state = O_STATE_WAIT;
			DBG("note on %d (%d)\r\n", args->note, s->ticks);
			seq_note_on(s, args);
		}
		m->duration -= 1;
		if (m->duration == 0) {
			// done
			m->op_state = O_STATE_INIT;
			DBG("note off (%d)\r\n", s->ticks);
			seq_note_off(s, args);
			return sizeof(struct note_args);
		}
	*/

	// waiting...
	return 0
}

// rest (op, duration)
func (m *seqModule) opRest(sm *seqStateMachine) int {
	/*
		struct rest_args *args = (struct rest_args *)&m->prog[m->pc];
		if (m->op_state == O_STATE_INIT) {
			// init
			m->duration = args->dur;
			m->op_state = O_STATE_WAIT;
		}
		m->duration -= 1;
		if (m->duration == 0) {
			// done
			m->op_state = O_STATE_INIT;
			return sizeof(struct rest_args);
		}

	*/
	// waiting...
	return 0
}

func (m *seqModule) tick(sm *seqStateMachine) {
	if sm.sstate == seqStateRun {
		opcode := sm.prog[sm.pc]
		switch opcode {
		case opNOP:
			sm.pc += m.opNOP(sm)
		case opLoop:
			sm.pc += m.opLoop(sm)
		case opNote:
			sm.pc += m.opNote(sm)
		case opRest:
			sm.pc += m.opRest(sm)
		default:
			panic(fmt.Sprintf("bad sequencer opcode %d", opcode))
		}
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
	prog     []uint8  // program memory
	pc       int      // program counter
	sstate   seqState // sequencer state
	ostate   opState  // operation state
	duration int      // operation duration
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
