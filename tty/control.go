/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   goAsk
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *  GoAsk ANSI Terminal gadgets.
 *-----------------------------------------------------------------*/
package tty

import "fmt"

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	// Progress Bar indicator styles
	ProgressStyle1 ProgressStyle = iota // rotating Braille characters
	ProgressStyle2                      // wave-style Braille characters
	ProgressStyle3                      // rotating slashdot
)

const (
	// Colored console ANSI colors control codes
	ansi_RESET         AnsiCode = "\033[0m"
	ansi_COLOR_RESET   AnsiCode = "\033[0m"
	ansi_RED           AnsiCode = "\033[31m"
	ansi_BRIGHT_RED    AnsiCode = "\033[91m"
	ansi_GREEN         AnsiCode = "\033[32m"
	ansi_BRIGHT_GREEN  AnsiCode = "\033[92m"
	ansi_YELLOW        AnsiCode = "\033[33m"
	ansi_BRIGHT_YELLOW AnsiCode = "\033[93m"
	ansi_PURPLE        AnsiCode = "\u001b[35m"
	ansi_BRIGHT_PURPLE AnsiCode = "\u001b[95m"
	ansi_CYAN          AnsiCode = "\u001b[36m"
	ansi_BRIGHT_CYAN   AnsiCode = "\u001b[96m"
	ansi_WHITE         AnsiCode = "\033[37m"
	ansi_BRIGHT_WHITE  AnsiCode = "\033[97m"
	ansi_HIDE_CURSOR   AnsiCode = "\033[?25l" // does not appear to work
	ansi_SHOW_CURSOR   AnsiCode = "\033[?25h" // idem
)

var (
	progressIndicators1 []rune = []rune{'⠇', '⠋', '⠉', '⠙', '⠸', '⠴', '⠤', '⠦'}
	progressIndicators2 []rune = []rune{'⠄', '⠤', '⠴', '⠶', '⠾', '⠿', '⠷', '⠶', '⠦', '⠤', '⠠'}
	progressIndicators3 []rune = []rune{'/', '-', '\\', '|'}
)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

// alias for ANSI control codes
type AnsiCode string

// alias for Progress Bar style
type ProgressStyle uint8

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

// Reset the terminal
func Reset() {
	fmt.Print("\033[39m\\033[49m")
}

// Move cursor to home position (top left)
func Home() {
	fmt.Print("\033[H")
}

// clear the screen but stay at cursor position
func ClearStay() {
	fmt.Print("\033[2K")
}

// clear the screen and move cursor to (0,0)
func Clear() {
	fmt.Print("\033[2J")
}

// clear everything from the cursor down
func ClearBelow() {
	fmt.Print("\033[0J")
}

// clear everything from the cursor up
func ClearAbove() {
	fmt.Print("\033[1J")
}

// Puts the cursor at that position but everything below is erased (?)
func Cursor(row, col int) {
	if row <= 0 {
		row = 1
	}
	if col <= 0 {
		col = 1
	}
	fmt.Printf("\033[%d;%dH", row, col)
}

// move the cursor N rows up
func CursorUp(n int) {
	fmt.Printf("\033[%dA", n)
}

// move the cursor N rows down
func CursorDown(n int) {
	fmt.Printf("\033[%dB", n)
}

// move the cursor N columns to the right
func CursorRight(n int) {
	fmt.Printf("\033[%dC", n)
}

// move the cursor N columns to the left
func CursorLeft(n int) {
	fmt.Printf("\033[%dD", n)
}

// Erase until the End-of-line
func EraseEOL() {
	fmt.Print("\033[K")
}

// Save cursor position
func SaveCursor() {
	fmt.Print("\033[s")
}

// Restore cursor position
func RestoreCursor() {
	fmt.Print("\033[u")
}

// print underlined text (no CR) and reset underline
func Underlined(str string) {
	fmt.Printf("\033[4m%s\033[24m", str)
}

// print bold text (no CR) and reset bold
func Bolded(str string) {
	fmt.Printf("\033[1m%s\033[21m", str)
}

// start using Bold
func Bold() {
	fmt.Print("\033[1m")
}

// terminate using Bold
func BoldOff() {
	fmt.Print("\033[22m")
}

// Print the args in color
func Color(color AnsiCode, args ...any) {
	fmt.Print(color)
	fmt.Print(args...)
	fmt.Print(ansi_COLOR_RESET)
}

// print in Red
func Red(args ...any) {
	Color(ansi_RED, args...)
}

// print in Green
func Green(args ...any) {
	Color(ansi_GREEN, args...)
}

// print in Yellow
func Yellow(args ...any) {
	Color(ansi_YELLOW, args...)
}

// Print in Magenta
func Purple(args ...any) {
	Color(ansi_PURPLE, args...)
}

// Print in Cyan
func Cyan(args ...any) {
	Color(ansi_CYAN, args...)
}

func BrightRed(args ...any) {
	Color(ansi_BRIGHT_RED, args...)
}

func BrightGreen(args ...any) {
	Color(ansi_BRIGHT_GREEN, args...)
}

func BrightYellow(args ...any) {
	Color(ansi_BRIGHT_YELLOW, args...)
}

func BrightPurple(args ...any) {
	Color(ansi_BRIGHT_PURPLE, args...)
}

func BrightCyan(args ...any) {
	Color(ansi_BRIGHT_CYAN, args...)
}

func BrightWhite(args ...any) {
	Color(ansi_BRIGHT_WHITE, args...)
}

// Show a progress bar of style at row of the console showing a title and
// an animation followed by the title. It is padded with "OK" when done.
func ShowProgressAt(style ProgressStyle, row int, title string, percent uint) {
	var pind []rune
	if style != ProgressStyle1 && style != ProgressStyle2 && style != ProgressStyle3 {
		style = ProgressStyle1
	}

	switch style {
	case ProgressStyle1:
		pind = progressIndicators1
	case ProgressStyle2:
		pind = progressIndicators2
	case ProgressStyle3:
		pind = progressIndicators3
	}

	var done string = ""
	if percent > 100 {
		println("percent > 100%")
		percent = 100
		return
	} else if percent == 100 {
		done = "✓ OK"
	}

	indicator := pind[int(percent)%len(pind)]
	fmt.Print(ansi_HIDE_CURSOR)
	Cursor(row, 0)
	fmt.Printf("%c %02d%%/100%% %4s %s\n", indicator, percent, done, title)
	fmt.Print(ansi_SHOW_CURSOR)
}
