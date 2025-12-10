/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   goAsk
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Simple Finite State Machine implementation for the purpose of
 * creating CLI interactive questionaires.
 *-----------------------------------------------------------------*/
package fsm

import "fmt"

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	StateNone  StateId = StateId(0)
	StateFinal StateId = StateId(65535) // pre-defined

/*
// These are user-provided in his/her module
InitialState StateId = iota + 1
State0
State1
FinalState
*/
)

var DefaultFinalState *State = NewState(StateFinal, "FIN", func(im IStateMachine) { fmt.Println("⚡ Entered FinalState") }, func(im IStateMachine) { fmt.Println("⚡ Exit FinalState") }, true, func(im IStateMachine) StateId { fmt.Println("⛔ The End"); return StateFinal })

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

// The State identifier type that is used when defining States. The module user
// should have a constant IOTA block that uses this type to enumerate states.
type StateId uint

// Callback function signature for a state's OnEnter event. When a State
// has an OnEnter event, it is always executed when the state is transitioned
// into from a previous state. It is skipped if the last run was the same
// state.
type OnEnterHandler func(IStateMachine)

// Callback function signature for a state's OnExit event.
type OnExitHandler func(IStateMachine)

// Callback function signature for a state's body which is executed
// in between OnEnter and OnExit.
type StateMainHandler func(IStateMachine) StateId

// An object representing a finite state machine's State.
type State struct {
	Id         StateId          // unique state identifier
	Name       string           // friendly name for the state
	parent     IStateMachine    // parent state machine
	onEnter    OnEnterHandler   // always executed prior to body on every State.Run()
	body       StateMainHandler // the main logic of the State
	onExit     OnExitHandler    // executed after Body but ONLY if there is a state transition
	isTerminal bool             // true if this is a terminal (end) state
}

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// (ctor) Creates a new instance of a State.
func NewState(id StateId, name string, onEnter OnEnterHandler, onExit OnExitHandler, terminal bool, body StateMainHandler) *State {
	return &State{
		Id:         id,
		Name:       name,
		parent:     nil,
		onEnter:    onEnter,
		body:       body,
		onExit:     onExit,
		isTerminal: terminal,
	}
}

// (ctor) simplified constructor for a simple state which only has a body but
// no OnEnter nor OnExit callbacks.
func NewStateSimple(id StateId, name string, terminal bool, body StateMainHandler) *State {
	return NewState(id, name, nil, nil, terminal, body)
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// implements fmt.Stringer
func (s *State) String() string {
	return s.Name
}

// executes a state
func (s *State) Run() StateId {
	// OnEnter is only executed upon the first transition
	if s.parent.GetPrevious() != s.Id && s.onEnter != nil {
		s.onEnter(s.parent)
	}

	// always execute the body of the state
	nextState := s.Id // self
	if s.body != nil {
		nextState = s.body(s.parent)
	}

	// OnExit is only executed if the FSM is transitioning
	// out of this state to another state. Never executed
	// for transitions to self
	if s.Id != nextState && s.onExit != nil {
		s.onExit(s.parent)
	}

	return nextState
}
