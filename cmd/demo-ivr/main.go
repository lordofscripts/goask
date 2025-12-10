/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   goAsk Demo #1
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Demonstration of various functionalities of the plain ASK module.
 * This demo uses the simple Finite State Machine to emulate an
 * Automated/Interactive Voice Response menu. In that system each
 * state makes use of goAsk to ask questions and transition to other
 * states depending on the answer.
 *  This is how you can use goAsk in a CLI application to interactively
 * ask the user and fill-in the CLI flags.
 *-----------------------------------------------------------------*/
package main

import (
	"fmt"
	"strconv"

	"github.com/lordofscripts/goask/ask"
	"github.com/lordofscripts/goask/fsm"
)

const (
	// These are user-provided in his/her module
	InitialState fsm.StateId = iota + 1
	BuyData
	BuyVoice
	TechSupport
	ClaimsDept
	Vacation
	FinalState
)

const (
	UNICODE_PHONE = rune(0x260e) // ☎
	UNICODE_MUSIC = rune(0x266c) // ♬
)

var _ fsm.OnEnterHandler = greeter
var _ fsm.OnExitHandler = byer

var myStateData MyUserData = MyUserData{Balance: 0.0, Taxes: 0.0, Fees: 0.0}

// implements IUserData
type MyUserData struct {
	Balance float32
	Taxes   float32
	Fees    float32
}

func txVoice(msg string, args ...any) {
	if len(args) == 0 {
		fmt.Printf("%c %s\n", UNICODE_PHONE, msg)
	} else {
		fmt.Printf("%c ", UNICODE_PHONE)
		fmt.Printf(msg, args...)
	}
}

func txMusic(msg string) {
	fmt.Printf("%c %s\n", UNICODE_MUSIC, msg)
}

func greeter(sm fsm.IStateMachine) {
	fmt.Println(rune(0x270b), " Hello!", sm.String())
}

func byer(sm fsm.IStateMachine) {
	fmt.Println(rune(0x270c), " Bye!", sm.String())
}

func getStateData(sm fsm.IStateMachine) *MyUserData {
	if v, ok := sm.GetStateData().(*MyUserData); ok {
		return v
	}
	panic("couldn't get state data!")
}

func defineStates() *fsm.StateMachine[MyUserData] {
	var sequencer *fsm.StateMachine[MyUserData]

	st0 := fsm.NewState(InitialState, "S0", greeter, nil, false, func(sm fsm.IStateMachine) fsm.StateId {
		fmt.Println(" * * * MAIN AUTOMATED RESPONSE MENU * * *")
		q1 := ask.NewMultipleChoiceQuestion("What would you like to do?", []ask.InputSelection{
			ask.NewInputSelection(0, "Hang up"),
			ask.NewInputSelection(1, "Buy Data packages"),
			ask.NewInputSelection(2, "Buy Voice packages"),
			ask.NewInputSelection(3, "Technical Support"),
			ask.NewInputSelection(4, "Claims Department"),
		})
		nextStateId := []fsm.StateId{FinalState, BuyData, BuyVoice, TechSupport, ClaimsDept}
		return nextStateId[q1.Ask().AsInt()]
	})

	calculateTaxAndFees := func(sm fsm.IStateMachine) {
		stateData := getStateData(sm)
		if stateData.Balance > 0 {
			stateData.Taxes = stateData.Balance * 0.14
			stateData.Fees = 15.0
		}
	}
	stX := fsm.NewState(FinalState, "SX", nil, byer, true, func(sm fsm.IStateMachine) fsm.StateId {
		if sm.GetPrevious() == InitialState {
			txVoice("Thank you for not bothering us more.\n")
		}
		calculateTaxAndFees(sm)
		data := getStateData(sm)
		if data.Balance > 0 {
			txVoice("You spent:\n\tPurchases: $%s\n\tTaxes    : $%s\n\tFees     : $%s\n\tTotal    : $%s\n",
				strconv.FormatFloat(float64(data.Balance), 'f', -1, 32),
				strconv.FormatFloat(float64(data.Taxes), 'f', -1, 32),
				strconv.FormatFloat(float64(data.Fees), 'f', -1, 32),
				strconv.FormatFloat(float64(data.Balance)+float64(data.Taxes)+float64(data.Fees), 'f', -1, 32))
		}
		txVoice("Thank you for choosing us!")
		fmt.Println("Hanging up...")
		return FinalState
	})

	st1 := fsm.NewState(BuyData, "SD", nil, nil, false, func(sm fsm.IStateMachine) fsm.StateId {
		q1 := ask.NewMultipleChoiceQuestion("Buy which DATA package?", []ask.InputSelection{
			ask.NewInputSelection(0, "Cancel"),
			ask.NewInputSelection(1, "3 days for $5"),
			ask.NewInputSelection(2, "7 days for $6"),
			ask.NewInputSelection(3, "10 days for $9"),
			ask.NewInputSelection(4, "Want VOICE instead"),
		})
		choice := q1.Ask().AsInt()
		var nextState fsm.StateId
		switch choice {
		case 0:
			nextState = InitialState
		case 4:
			nextState = BuyVoice
		default:
			costs := []float32{0, 5.0, 6.0, 9.0}
			stateData := getStateData(sm)
			stateData.Balance = costs[choice]
			nextState = FinalState
		}
		return nextState
	})

	st2 := fsm.NewState(BuyVoice, "SV", nil, nil, false, func(sm fsm.IStateMachine) fsm.StateId {
		q1 := ask.NewMultipleChoiceQuestion("Buy which VOICE package?", []ask.InputSelection{
			ask.NewInputSelection(0, "Cancel"),
			ask.NewInputSelection(1, "1 week for $7"),
			ask.NewInputSelection(2, "2 weeks for $14"),
			ask.NewInputSelection(3, "1 month for $35"),
			ask.NewInputSelection(4, "Want DATA instead"),
		})
		choice := q1.Ask().AsInt()
		var nextState fsm.StateId
		switch choice {
		case 0:
			nextState = InitialState
		case 4:
			nextState = BuyData
		default:
			costs := []float32{0, 7.0, 14.0, 35.0}
			stateData := getStateData(sm)
			stateData.Balance = costs[choice]
			nextState = FinalState
		}
		return nextState
	})

	st3 := fsm.NewState(TechSupport, "ST", nil, nil, false, func(sm fsm.IStateMachine) fsm.StateId {
		q1 := ask.NewMultipleChoiceQuestion("Select Technical Support area?", []ask.InputSelection{
			ask.NewInputSelection(0, "Cancel"),
			ask.NewInputSelection(1, "Internet"),
			ask.NewInputSelection(2, "Phone"),
			ask.NewInputSelection(3, "Cable TV"),
		})
		choice := q1.Ask().AsInt()
		var nextState fsm.StateId
		switch choice {
		case 0:
			nextState = InitialState
		default:
			nextState = Vacation
		}
		return nextState
	})

	st4 := fsm.NewState(ClaimsDept, "SC", nil, nil, false, func(sm fsm.IStateMachine) fsm.StateId {
		q1 := ask.NewMultipleChoiceQuestion("Which type of claim?", []ask.InputSelection{
			ask.NewInputSelection(0, "Cancel"),
			ask.NewInputSelection(1, "Sales returns"),
			ask.NewInputSelection(2, "Customer service complaints"),
			ask.NewInputSelection(3, "General feedback"),
		})
		choice := q1.Ask().AsInt()
		var nextState fsm.StateId
		switch choice {
		case 0:
			nextState = InitialState
		default:
			nextState = Vacation
		}
		return nextState
	})

	st5 := fsm.NewState(Vacation, "SF", func(im fsm.IStateMachine) {
		txVoice("The department you have reached is on vacation. Call another time.")
	}, func(im fsm.IStateMachine) {
		txVoice("We are transferring you to another department. Trust us!")
	}, false, func(im fsm.IStateMachine) fsm.StateId {
		txMusic("(Annoying music here)")
		return fsm.StateFinal // use the default
	})
	// we have two final states, a nice one and a rude one
	sequencer = fsm.NewStateMachine[MyUserData]("Customer Service", st0, st1, st2, st3, st4, st5, stX, fsm.DefaultFinalState).SetUserDataObject(&myStateData)

	return sequencer
}

func main() {
	sm := defineStates()
	sm.Start()
}
