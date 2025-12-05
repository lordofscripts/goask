/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   go-ask
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *  GoAsk is a simple module made with command-line applications in
 * mind. It has multiple objects representing questions, answers and
 * questionaires.
 *  Rather than relying on a plain CLI which can get confusing to
 * inexperienced or unfamiliar users, goAsk allows the application
 * programmer to design a decision path of questions where the answers
 * are mapped to the CLI flags. So, with this you can have an interactive
 * dialog with the user to setup CLI options.
 *-----------------------------------------------------------------*/
package goask

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	// Colored console ANSI colors control codes
	ANSI_RESET  string = "\033[0m"
	ANSI_RED    string = "\033[31m"
	ANSI_GREEN  string = "\033[32m"
	ANSI_YELLOW string = "\033[93m"
	ANSI_PURPLE string = "\u001b[35m"
	ANSI_BROWN  string = "\u001b[33m"

	// Unicode symbols
	ICON_WHITE_RIGHT rune = rune(0x1f449) // 👉
)
