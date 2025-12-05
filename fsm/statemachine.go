/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   goAsk
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Simple Finite State Machine implementation suitable for questionaires.
 *-----------------------------------------------------------------*/
package fsm

import (
	"fmt"
	"sync"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ IStateMachine = (*StateMachine[any])(nil)

// Interface of a simple Finite State Machine
type IStateMachine interface {
	fmt.Stringer
	// get the friendly name of the FSM
	GetName() string
	// get the previous state ID.
	GetPrevious() StateId
	// Whether the FSM is running
	IsActive() bool
	// Whether the FSM reached its end-of-life
	IsDone() bool
	// get the state data. The caller must cast it to the proper type
	GetStateData() any
}

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

// An object representing a simple Finite State Machine that has
// just enough features to implement interactive questionaires.
type StateMachine[T any] struct {
	name          string
	states        map[StateId]*State
	isActive      bool
	isFinished    bool
	initialState  *State
	previousState StateId
	stateData     *T
	mu            sync.Mutex
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// (ctor) creates a new instance of a Finite State Machine with an
// initial state and the other states.
func NewStateMachine[T any](name string, initialState *State, otherStates ...*State) *StateMachine[T] {
	me := &StateMachine[T]{
		name:          name,
		states:        make(map[StateId]*State),
		isActive:      false,
		isFinished:    false,
		initialState:  initialState,
		previousState: StateNone,
		stateData:     nil,
	}

	// compose the list of states
	initialState.parent = me
	me.states[initialState.Id] = initialState
	for _, state := range otherStates {
		state.parent = me // assign parent FSM (this instance)
		me.states[state.Id] = state
	}

	return me
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// implements fmt.Stringer
func (sm *StateMachine[T]) String() string {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	return fmt.Sprintf("[%s] active:%t done:%t", sm.name, sm.isActive, sm.isFinished)
}

// integrity check that the defined state machine is good. Must be checked
// prior to Start()
func (sm *StateMachine[T]) IsValid() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if len(sm.states) == 0 {
		return fmt.Errorf("no states have been defined")
	}
	if len(sm.states) == 1 {
		return fmt.Errorf("only has an initial state")
	}

	for _, state := range sm.states {
		if state.isTerminal {
			return nil
		}
	}

	return fmt.Errorf("there must be at least one terminal state")
}

// get the FSM's friendly name
func (sm *StateMachine[T]) GetName() string {
	return sm.name
}

// get the previous StateId. It can be used within State.OnEnter
// to check what the previous state was.
func (sm *StateMachine[T]) GetPrevious() StateId {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	return sm.previousState
}

// Whether the FSM is running
func (sm *StateMachine[T]) IsActive() bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	return sm.isActive
}

// Whether the FSM reached its end-of-life
func (sm *StateMachine[T]) IsDone() bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	return sm.isFinished
}

// get the custom state data that could be checked and modified
// during State.Run(). Can be set with SetUserDataObject()
func (sm *StateMachine[T]) GetStateData() any {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	return sm.stateData
}

// Set the custom state data
func (sm *StateMachine[T]) SetUserDataObject(userData *T) *StateMachine[T] {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.stateData = userData
	return sm
}

// Start executing the State machine
func (sm *StateMachine[T]) Start() {
	sm.isActive = true
	sm.isFinished = false

	// the initial state is always executed first
	nextState := sm.initialState.Run()
	// update because InitialState previous's is StateNone
	sm.previousState = sm.initialState.Id

	// Loop through the defined states transitions
	for !sm.isFinished {
		currentState := sm.states[nextState]
		nextState = currentState.Run()
		// update the previous state if we are transitioning
		if nextState != sm.previousState {
			sm.previousState = currentState.Id
		}
		// check if terminating by FSM definition
		sm.isFinished = currentState.isTerminal
	}

	sm.isActive = false
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/
