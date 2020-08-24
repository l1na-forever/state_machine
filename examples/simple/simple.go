package main

import (
	sm "github.com/l1na-forever/state_machine"
	"log"
	"time"
)

const (
	RunningStateTime = time.Duration(5) * time.Second
)

// Left here as global for simplicity of the example
var (
	machine      *sm.StateMachine
	idleState    = &IdleState{}
	runningState = &RunningState{}
)

type IdleState struct{}

func (i *IdleState) EnterAllowed(sm.State) bool { return true }
func (i *IdleState) ExitAllowed(sm.State) bool  { return true }
func (i *IdleState) Enter()                     { log.Println("Entered idle") }
func (i *IdleState) Exit()                      { log.Println("Exited idle") }

type RunningState struct {
	entered time.Time
}

func (r *RunningState) EnterAllowed(sm.State) bool { return r.runCompleted() }
func (r *RunningState) ExitAllowed(sm.State) bool  { return r.runCompleted() }

func (r *RunningState) Enter() {
	r.entered = time.Now()
	time.AfterFunc(RunningStateTime, func() { machine.Transition(idleState) })
	log.Println("Started running")
}

func (r *RunningState) Exit() {
	log.Println("Done running")
}

func (r *RunningState) runCompleted() bool {
	// State runs for 5 seconds after being entered
	return time.Since(r.entered) >= RunningStateTime
}

func main() {
	machine, err := sm.NewStateMachine(idleState)
	if err != nil {
		log.Fatal("Couldn't set up state machine", err)
	}

	// Transition whenever needed:
	err = machine.Transition(runningState)
	if err != nil {
		log.Println("Couldn't transition to running", err)
	}
}
