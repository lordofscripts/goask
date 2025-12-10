/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   goAsk Demo #1
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Demonstration of various functionalities of the plain ASK module.
 * Things such as asking simple values, multiple choice, etc.
 *-----------------------------------------------------------------*/
package main

import (
	"fmt"
	"os"

	"github.com/lordofscripts/goask"
	"github.com/lordofscripts/goask/ask"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				I n i t i a l i z e r
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

func askPlainQuestions() {
	fmt.Println("*** Plain Questions ***")

	intInp := ask.NewIntInputRequest("Enter integer", 5)
	intInp.Read()
	fmt.Println("· Int input", intInp.Value)

	strInp := ask.NewStringInputRequest("Enter string", "test")
	strInp.Read()
	fmt.Println("· Str input", strInp.Value)

	runInp := ask.NewRuneInputRequest("Enter char", 'x')
	runInp.Read()
	fmt.Println("· Rune input", string(runInp.Value))
}

func askMultipleChoice() {
	fmt.Println("*** Multiple Choice (single) ***")

	question1 := ask.NewMultipleChoiceQuestion("Which one do you want?", []ask.InputSelection{
		ask.NewInputSelection(0, "Default"),
		ask.NewInputSelection(1, "One"),
		ask.NewInputSelection(2, "Two"),
	})
	question1.Ask()
	fmt.Println("· Question #1: ", question1.AsInt())

	wants := ask.SelectOptions("Which one you want?", []ask.InputSelection{
		{Number: 0, Text: "Default"},
		{Number: 1, Text: "One"},
		{Number: 2, Text: "Two"},
	})

	fmt.Println("· User selected", wants)
}

/* ----------------------------------------------------------------
 *						T e s t s
 *-----------------------------------------------------------------*/

func OpHelp() {
	fmt.Println(goask.ANSI_BROWN)
	fmt.Println("I am helping you")
	OpGoodBye()
}

func OpList() {
	fmt.Println(goask.ANSI_BROWN)
	fmt.Println("I am listing")
	OpGoodBye()
}

func OpGoodBye() {
	fmt.Println("Good bye!")
	fmt.Println(goask.ANSI_RESET)
	os.Exit(0)
}

func OpEncrypt(alpha, supplemental, cipher, text string) {
	fmt.Println("Encrypt")
	fmt.Printf("\tAlpha: %s\n\tCipher: %s\n\tText: %s\n", alpha, cipher, text)
}

func OpDecrypt(alpha, supplemental, cipher, text string) {
	fmt.Println("Decrypt")
	fmt.Printf("\tAlpha: %s\n\tCipher: %s\n\tText: %s\n", alpha, cipher, text)
}

/* ----------------------------------------------------------------
 *					M A I N    |     D E M O
 *-----------------------------------------------------------------*/

func main() {
	askPlainQuestions()
	askMultipleChoice()
}
