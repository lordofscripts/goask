/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package ask

import (
	"bufio"
	"fmt"
	"lordofscripts/goask"
	"os"
	"slices"
	"strconv"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ ICurious = (*QuestionWithChoice)(nil)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type QuestionWithChoice struct {
	Prompt  string
	Choices []InputSelection
	answer  int
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// (ctor) creates a multiple choice question
func NewMultipleChoiceQuestion(prompt string, choices []InputSelection) *QuestionWithChoice {
	return &QuestionWithChoice{
		Prompt:  prompt,
		Choices: choices,
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// implements ask.ICurious and returns the answer.
func (q *QuestionWithChoice) Answer() any {
	return q.answer
}

// implements ask.ICurious and uses stdin (console) to ask the
// user to select a valid choice. It keeps on asking until a
// valid option is chosen.
func (q *QuestionWithChoice) Ask() ICurious {
	q.answer = -1
	if len(q.Choices) == 0 {
		return q
	}
	if len(q.Choices) == 1 {
		q.answer = int(q.Choices[0].Number)
		return q
	}

	// list of valid option numbers
	var valid []uint = make([]uint, 0)
	for _, opt := range q.Choices {
		if !slices.Contains(valid, opt.Number) {
			valid = append(valid, uint(opt.Number))
		}
	}

	renderMenu := func() {
		fmt.Println(goask.ANSI_YELLOW, q.Prompt, goask.ANSI_GREEN)
		for i, opt := range q.Choices {
			var isDef string = ""
			if i == 0 {
				isDef = "(default)"
			}
			fmt.Printf("\t%d. %s %s\n", opt.Number, opt.Text, isDef)
		}
		fmt.Print(goask.ANSI_RESET)
	}

	readSelection := func() int {
		fmt.Print("Enter your choice: ")
		reader := bufio.NewReader(os.Stdin)
		if str, err := reader.ReadString('\n'); err == nil {
			str = str[:len(str)-1]
			if len(str) == 0 {
				return int(q.Choices[0].Number)
			} else if nr, err := strconv.Atoi(str); err != nil {
				return -1
			} else {
				return nr
			}
		}
		return -1
	}

	selected := -1
	for selected == -1 {
		renderMenu()
		if value := readSelection(); value > -1 && slices.Contains(valid, uint(value)) {
			selected = value
		}
	}

	q.answer = selected
	fmt.Println(q.Choices[selected].Chosen())
	return q
}

// implements ask.ICurious and returns the value
// as the chosen option number
func (q *QuestionWithChoice) AsInt() int {
	return q.answer
}

// implements ask.ICurious and returns the value
// as a rune
func (q *QuestionWithChoice) AsRune() rune {
	return rune(strconv.Itoa(q.answer)[0])
}

// implements ask.ICurious and returns the text
// of the chosen answer
func (q *QuestionWithChoice) AsString() string {
	return q.Choices[q.answer].Text
}
