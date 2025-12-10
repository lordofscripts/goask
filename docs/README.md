# Go Ask !!!

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/lordofscripts/goask)
[![Go Report Card](https://goreportcard.com/badge/github.com/lordofscripts/goask?style=flat-square)](https://goreportcard.com/report/github.com/lordofscripts/goask)
[![Go Reference](https://pkg.go.dev/badge/github.com/lordofscripts/goask.svg)](https://pkg.go.dev/github.com/lordofscripts/goask)
[![GitHub release (with filter)](https://img.shields.io/github/v/release/lordofscripts/goask)](https://github.com/lordofscripts/goask/releases/latest)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL_2.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)

![](./gopher-goask.png)

Go ask came to life as a necessity to another GO command-line application
in which the combination of command-line flags could be a bit too much
for users not accustomed to the CLI. Therefore, it occured to me that it
could be nice to have a CLI interactive questionaire that would interview
the user, and as a result fill-in the corresponding CLI flags accordingly.

## Install

To install into your project's `go.mod`:

> go get github.com/lordofscripts/goask

And to use in your GO code:

> import "github.com/lordofscripts/goask/ask"
> import "github.com/lordofscripts/goask/fsm"

## Ask Package

It contains the main functionality where the Go module asks, via the
console, the user for values.

### Plain inputs

For simple input requests use the `ask.InputRequest` type. It can ask the
user via stdin (console) any integer, string, or rune (single character).
Use the appropriate constructor for the input type and ask. Retrieve the
result using any of the `AsInt()`, `AsString()` or `AsRune()` depending
on the type you choose to ask.

> request := ask.NewIntInputRequest("Enter value", 0)
> value := request.Ask().AsInt()
> fmt.Printf("The answer was %d (%d)\n", value, request.Answer())

One way to construct a multiple-choice question is via the `ask.InputSelection`
type. It lets you define its numerical id/value and its corresponding text
to be displayed. Then you can use the 

> options := []InputSelection{
>   ask.NewInputSelection(0, "Cancel"),
>   ask.NewInputSelection1, "Colored"),
>   ask.NewInputSelection2, "Black & White"),
> }
> nr := ask.SelectOptions("Please choose", options)
> fmt.Printf("You selected #%d: %s\n", nr, options[nr].Chosen())

### Multiple choice questions

Alternatively, and specially if you are going to use it as a `SmartQuestion`
you can use `ask.QuestionWithChoice`. So using the same options as above:

> mchoice := ask.NewMultipleChoiceQuestion("Please choose", options)
> value := mchoice.Ask().AsInt()
> fmt.Printf("You selected #%d: %s\n", mchoice.Answer(), options[nr].Chosen())

## Questionaires

When you have questions the logical follow up would be a questionaire. This
is quite useful in CLI applications when the amount of CLI flags or their
combinations can get to overwhelming for non-experienced users (Windowers).
I came up with the idea of using a simple Finite State Machine (FSM) to
implement questionaires.

With a simple FSM you could define several `fsm.State` of a FSM 
`fsm.StateMachine[T]`. Depending on the answers given by the user on any
given `State` the FSM transitions to another state. In this way you can
implement both simple and complex questionaires. In the state body function
you would use the functionality of the `ask` package as you can see in
the `cmd/demo-ivr` sample application that creates a typical Interactive
Voice Response emulator.

### Finite State Machine

Organize your flow of questions and answers into **states**. Enumerate each
of the states by using IOTA starting with `1` as value because `fsm.StateNone`
is predefined.

> const (
>   InitialState fsm.StateId = iota + 1
>   State1
>   State2
>   FinalState     
> )

In every state you may use *state data* that could be shared among states.
If not simply define an empty structure.

> type myStateData struct {}

Then declare the `StateMachine[T]` without details yet because we would
need to pass this to the constructors of `State` instances:

> // The instances implement the fsm.IStateMachine interface
> var fsm *fsm.StateMachine[myUserData]

Then define the initial state because you will need to supply this 
separately when instantiating the State Machine object. Subsequently
define the other states accordingly. Keep in mind that at least one
of the states (that is not the initial state) has to be marked as final
or terminal.

A terminal state is where the State Machine terminates execution or
its useful life. A State Machine has at least ONE terminal state. But
depending on the logic of your system, there may be multiple terminal
states.

Now that you created each of your states with their corresponding logic,
you would instantiate the State machine as follows (remember we 
previously declared it):

> fsm = fsm.NewStateMachine[myUserData]("Customer Service", 
>           stInitial, st1, st2, stFinal).
>       SetUserDataObject(&myStateData)
 
And start running it synchronously:

> fsm.Start()

Because this FSM is meant for the purpose of questionaires, it does not
implement fancy features of a proper FSM like asynchronous execution,
events, or locking.

#### Creating the State

A state has an ID, a friendly name, whether it is terminal or not, and
several callback functions that define what it does.

> type OnEnterHandler func(IStateMachine)

Every state can have an `onEnter` optional event handler. It is executed only
when the state is entered after transitioning from a different state.
Therefore, it is not executed with *transitions to self*. Else don't
use that and specify it as `nil` in the `State` constructor.

> type OnExitHandler func(IStateMachine)

Likewise, every state has an `onExit` optional handler. The handler
is executed only when the state is transitioning to a different state.

> type StateMainHandler func(IStateMachine) StateId

And every state **must** have a body executed after the enter event and
before the exit event. This body contains the main decision logic of
the state. This state analyzes the input, in GoAsk the input is the
answers to the questions, and depending on the answer it decides to
which next state it will transition on the next iteration. That is
what the return value is.

This is how you could construct the state, in this example the
initial state. You can specify the handlers via pointers to functions
or inline. For the sake of clarity I will use the former.

> func initialStateEnter(sm fsm.IStateMachine) {
>   fmt.Print("Starting up!")   
> }
>
> func initialStateExit(sm fsm.IStateMachine) {
>   fmt.Print("Finished warming up!")   
> }
>
> func initialStateBody(sm fsm.IStateMachine) fsm.StateId {
> 	fmt.Println(" * * * MAIN AUTOMATED RESPONSE MENU * * *")
>   q1 := ask.NewMultipleChoiceQuestion("What would you like to do?", []ask.InputSelection{
>       ask.NewInputSelection(0, "Hang up"),
>       ask.NewInputSelection(1, "Buy Data packages"),
>       ask.NewInputSelection(2, "Buy Voice packages"),
>       ask.NewInputSelection(3, "Technical Support"),
>       ask.NewInputSelection(4, "Claims Department"),
>   })
>   nextStateId := []fsm.StateId{FinalState, BuyData, BuyVoice, TechSupport, ClaimsDept}
>   return nextStateId[q1.Ask().AsInt()]
> }

And then create the state.

> stInitial := fsm.NewState(InitialState, "S0", 
>                   initialStateEnter, // optional
>                   initialStateExit,  // optional
>                   false, // true only for terminal states!
>                   initialStateBody

If you don't want to bother creating your own final state, you can use
a predefined one with Id `fsm.StateFinal` that has already been
instantiated as `fsm.DefaultFinalState`.