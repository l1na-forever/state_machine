state_machine
=============

Dead simple state machine for programs that mostly manage transitions themselves.

## Installation

    go get github.com/l1na-forever/state_machine

To run included tests:

    cd $GOROOT/src/pkg/github.com/l1na-forever/state_machine
    go test

## Status

Functional for my needs. state_machine is capable of:

  - Transitioning between state objects that conform to the State interface
  - Fitting its entire implementation on one screen

In addition,

  - state_machine is [fully documented](https://godoc.org/github.com/l1na-forever/state_machine)
  - state_machine has thorough test coverage

## Quick Start

Start by implementing the [State interface](https://godoc.org/github.com/l1na-forever/state_machine/#State) to set up your various states.

```go
import (
  "log"
  "time"
)

const (
  RunningStateTime = time.Duration(5) * time.Second
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
  log.Println("Started running")
}

func (r *RunningState) Exit() {
  log.Println("Done running")
}

func (r *RunningState) runCompleted() bool {
  // State runs for 5 seconds after being entered
  return time.Since(r.entered) >= RunningStateTime
}
```

Afterwards, the state machine itself can be initialized. Outside of the active state, state_machine maintains no references to the universe of states available in the system; the program should maintain these references itself (as will be needed to call [Transition](https://godoc.org/github.com/l1na-forever/state_machine/#StateMachine.Transition)).

```go
import (
  sm "github.com/l1na-forever/state_machine"
)

// Left here as global for simplicity of the example
var (
  machine      *sm.StateMachine
  idleState    = &IdleState{}
  runningState = &RunningState{}
)

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
```

With the full above example, the running state is entered, but never transitions back to the idle state. This is expected; state_machine won't initiate transitions on its own. To complete the example, the running state should signal its own transition back to idle.

```go
func (r *RunningState) Enter() {
  r.entered = time.Now()
  time.AfterFunc(RunningStateTime, func() { machine.Transition(idleState) })
  log.Info("Started running")
}
```

state_machine's usefulness comes into play when another area (or thread) of the program would have made an incorrect transition:

```go
  machine.Transition(runningState)
  // This second transition will fail, resulting in an error, as RunningState's ExitAllowed will return false.
  err = machine.Transition(idleState)
```

**Rather than having a static mapping of allowed transitions, States are given the opportunity to respond to each transition request with a distinct outcome.**

## Documentation

API Documentation can be found at [godoc.org/github.com/l1na-forever/state_machine](https://godoc.org/github.com/l1na-forever/state_machine).

## License

Copyright © 2020 Lina

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

