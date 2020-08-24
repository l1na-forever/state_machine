package state_machine

import (
	"testing"
)

var _ State = &EmptyState{}

func TestNewStateMachineNilState(t *testing.T) {
	_, err := NewStateMachine(nil)
	if err != ErrNilInitialState {
		t.Fail()
	}
}

func TestNewStateMachineUnusableState(t *testing.T) {
	_, err := NewStateMachine(&stateEnterNotAllowed{})
	if err != ErrUnusableInitialState {
		t.Fail()
	}
}

func TestNewStateMachineHappy(t *testing.T) {
	initialState := &stateAllAllowed{}
	machine, err := NewStateMachine(initialState)
	if err != nil {
		t.Fatal(err)
	}

	if machine == nil || machine.ActiveState() != initialState {
		t.Fail()
	}
	if !initialState.Entered {
		t.Fail()
	}
}

func TestTransitionExitNotAllowed(t *testing.T) {
	source := &stateExitNotAllowed{}
	destination := &stateAllAllowed{}
	machine, err := NewStateMachine(source)
	if err != nil {
		t.Fatal(err)
	}

	err = machine.Transition(destination)
	if err != ErrTransitionDenied {
		t.Fail()
	}
	if destination.Entered {
		t.Fail()
	}
}

func TestTransitionEnterNotAllowed(t *testing.T) {
	source := &stateAllAllowed{}
	destination := &stateEnterNotAllowed{}
	machine, err := NewStateMachine(source)
	if err != nil {
		t.Fatal(err)
	}

	err = machine.Transition(destination)
	if err != ErrTransitionDenied {
		t.Fail()
	}
	if source.Exited {
		t.Fail()
	}
}

func TestTransitionHappy(t *testing.T) {
	source := &stateAllAllowed{}
	destination := &stateAllAllowed{}
	machine, err := NewStateMachine(source)
	if err != nil {
		t.Fatal(err)
	}

	err = machine.Transition(destination)
	if err != nil {
		t.Fatal(err)
	}
	if machine.ActiveState() != destination {
		t.Fail()
	}
	if !source.Exited {
		t.Fail()
	}
	if !destination.Entered {
		t.Fail()
	}
}

type stateExitNotAllowed struct{}

func (s *stateExitNotAllowed) EnterAllowed(State) bool { return true }
func (s *stateExitNotAllowed) ExitAllowed(State) bool  { return false }
func (s *stateExitNotAllowed) Enter()                  {}
func (s *stateExitNotAllowed) Exit()                   { panic("") }

type stateEnterNotAllowed struct{}

func (s *stateEnterNotAllowed) EnterAllowed(State) bool { return false }
func (s *stateEnterNotAllowed) ExitAllowed(State) bool  { return true }
func (s *stateEnterNotAllowed) Enter()                  { panic("") }
func (s *stateEnterNotAllowed) Exit()                   {}

type stateAllAllowed struct {
	Entered, Exited bool
}

func (s *stateAllAllowed) EnterAllowed(State) bool { return true }
func (s *stateAllAllowed) ExitAllowed(State) bool  { return true }
func (s *stateAllAllowed) Enter()                  { s.Entered = true }
func (s *stateAllAllowed) Exit()                   { s.Exited = true }
