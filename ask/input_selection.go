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
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type InputSelection struct {
	Number uint
	Text   string
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// (ctor) a new (value-type) input selection.
func NewInputSelection(nr uint, text string) InputSelection {
	return InputSelection{
		Number: nr,
		Text:   text,
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// implements fmt.Stringer
// the input selection in the form of "Number. Text" suitable for
// rendering in a menu.
func (is InputSelection) String() string {
	return fmt.Sprintf("%d. %s", is.Number, is.Text)
}

// utility method to format a chosen option
func (is InputSelection) Chosen() string {
	return fmt.Sprintf("%c %s", goask.ICON_WHITE_RIGHT, is.Text)
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

// given the options, show the prompt and enumerate all the options.
// then use stdin to ask the user until a valid option number is
// selected.
func SelectOptions(prompt string, options []InputSelection) int {
	if len(options) == 0 {
		return -1
	}
	if len(options) == 1 {
		return int(options[0].Number)
	}

	// list of valid option numbers
	var valid []uint = make([]uint, 0)
	for _, opt := range options {
		if !slices.Contains(valid, opt.Number) {
			valid = append(valid, uint(opt.Number))
		}
	}

	renderMenu := func() {
		fmt.Println(goask.ANSI_YELLOW, prompt, goask.ANSI_GREEN)
		for i, opt := range options {
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
				return int(options[0].Number)
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

	fmt.Printf("%c %d\n", goask.ICON_WHITE_RIGHT, selected)
	return selected
}
