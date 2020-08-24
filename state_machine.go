// Dead simple state machine for programs that mostly manage transitions
// themselves.
package state_machine

import (
	"errors"
	"sync"
)

var (
	ErrTransitionDenied     = errors.New("Transition was denied by ExitAllowed/EnterAllowed")
	ErrNilInitialState      = errors.New("An initial state must be given")
	ErrUnusableInitialState = errors.New("The initial state did not allow an Enter transition")
)

// A program's states should conform to this interface to be used with
// state_machine. All callbacks are called synchronously.
type State interface {
	// Return true if the state machine may transition into this state. The
	// source (active) state is given as an argument.
	EnterAllowed(source State) bool

	// Return true if the state machine may transition away from this state. The
	// proposed destination state is given as an argument.
	ExitAllowed(destination State) bool

	// Called when this state is being transitioned into.
	Enter()

	// Called when this state is being transitioned away from.
	Exit()
}

// An empty default state that allows, and no-ops, all transitions. Useful for
// building less complex states that do not need a custom implementation for
// all callbacks.
type EmptyState struct{}

func (e *EmptyState) EnterAllowed(source State) bool { return true }
func (e *EmptyState) ExitAllowed(source State) bool  { return true }
func (e *EmptyState) Enter()                         {}
func (e *EmptyState) Exit()                          {}

// Construct with NewStateMachine only.
type StateMachine struct {
	activeState State
	mutex       sync.Mutex
}

// Construct a new state machine. An initial state must be supplied. The
// initial state must allow a transition to itself (returning 'true' from
// EnterAllowed), and the initial state's Enter callback will be called before
// the machine is constructed.
func NewStateMachine(initialState State) (*StateMachine, error) {
	if initialState == nil {
		return nil, ErrNilInitialState
	}

	if !initialState.EnterAllowed(initialState) {
		return nil, ErrUnusableInitialState
	}
	initialState.Enter()

	return &StateMachine{
		activeState: initialState,
	}, nil
}

// Returns the active machine state. A state is considered active only once
// its transition event has completed both the outgoing state's Exit callback
// and the incoming state's Enter callback.
func (s *StateMachine) ActiveState() State {
	return s.activeState
}

// Attempt to transition to the given destination state. If both states do not
// accept the transition (one of the states returned 'false' from their
// ExitAllowed/EnterAllowed callbacks), ErrTransitionDenied will be returned.
//
// The machine's active state is updated only once both the outgoing state and
// incoming state's Exit and Enter callbacks have synchronously completed.
func (s *StateMachine) Transition(destination State) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.activeState.ExitAllowed(destination) && destination.EnterAllowed(s.activeState) {
		s.activeState.Exit()
		destination.Enter()
		s.activeState = destination
	} else {
		return ErrTransitionDenied
	}

	return nil
}
