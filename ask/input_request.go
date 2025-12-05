/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   goAsk
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * An InputRequest implements ICurious and uses stdin to request a
 * value from the user. The value can be an integer, a rune or string.
 *-----------------------------------------------------------------*/
package ask

import (
	"bufio"
	"fmt"
	"lordofscripts/goask"
	"os"
	"strconv"
	"strings"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ ICurious = (*InputRequest[int])(nil)
var _ ICurious = (*InputRequest[rune])(nil)
var _ ICurious = (*InputRequest[string])(nil)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

// A generic constraint for the InputRequest
type NumberStringRune interface {
	int | string | rune
}

// an input request of the supported types (int,string,rune)
type InputRequest[T NumberStringRune] struct {
	Prompt  string
	Default T
	Value   T
}

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// (ctor) request an integer value
func NewIntInputRequest(prompt string, defval int) *InputRequest[int] {
	return &InputRequest[int]{prompt, defval, 0}
}

// (ctor) request a string value
func NewStringInputRequest(prompt string, defval string) *InputRequest[string] {
	return &InputRequest[string]{prompt, defval, ""}
}

// (ctor) request a rune value
func NewRuneInputRequest(prompt string, defval rune) *InputRequest[rune] {
	return &InputRequest[rune]{prompt, defval, rune(0)}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// ask for the value. To obtain the answer use any of Answer(),
// AsInt(), AsRune() or AsString() depending on the value type.
// to retrieve the value immediately use Read() instead.
func (r *InputRequest[T]) Ask() ICurious {
	r.Read()
	return r
}

// obtain the answer to the question
func (r *InputRequest[T]) Answer() any {
	return r.Value
}

// obtain the answer as an integer. Only valid if the
// answer type is integer, else it prints an error and
// returns -1.
func (r *InputRequest[T]) AsInt() int {
	if val, ok := any(r.Value).(int); ok {
		return val
	}
	println("InputRequest is not int value")
	return -1
}

// obtain the answer as a string. This works regardless
// of the answer's underlying type.
func (r *InputRequest[T]) AsString() string {
	var result string
	switch v := any(r.Value).(type) {
	case rune:
		result = string(v)

	case int:
		result = strconv.Itoa(v)

	case string:
		result = v

	default:
		result = ""
	}

	return result
}

// obtain the answer as a rune. If the answer's
// underlying value is not a rune, it returns 0
// after printing an error.
func (r *InputRequest[T]) AsRune() rune {
	if val, ok := any(r.Value).(rune); ok {
		return val
	}

	println("InputRequest is not rune value")
	return 0
}

// read the answer from stdin and return the answer
// to the caller. The same result can be retrieved later
// by calling Answer().
func (r *InputRequest[T]) Read() T {
	var err error = nil
	var value T
	switch v := any(r.Default).(type) {
	case int:
		requestInteger := func() error {
			var n int
			fmt.Printf("%s [%d]: ", r.Prompt, v)
			if n, err = fmt.Scanln(&value); err != nil {
				if n == 0 {
					value = r.Default
					err = nil
				} else {
					fmt.Printf("!!! Error reading input: %v\n", err)
				}

				return err
			}

			result, _ := any(value).(int)
			fmt.Printf("%c %d\n", goask.ICON_WHITE_RIGHT, result)
			return nil
		}
		err := requestInteger()
		for err != nil {
			err = requestInteger()
		}

	case string:
		fmt.Printf("%s [%s]: ", r.Prompt, v)
		var str string
		reader := bufio.NewReader(os.Stdin)
		str, err = reader.ReadString('\n')
		if len(strings.Trim(str, " \t\n")) == 0 {
			value = r.Default
		} else {
			str = str[:len(str)-1]
			value = any(str).(T)
		}
		result, _ := any(value).(string)
		fmt.Printf("%c %s\n", goask.ICON_WHITE_RIGHT, result)

	case rune:
		fmt.Printf("%s [%c]: ", r.Prompt, v)
		reader := bufio.NewReader(os.Stdin)
		str, err := reader.ReadString('\n')
		if len(strings.Trim(str, " \t\n")) == 0 {
			value = r.Default
		} else {
			str = str[:len(str)-1]
			if err == nil && len(str) > 0 {
				value = any([]rune(str)[0]).(T)
			}
		}
		result, _ := any(value).(rune)
		fmt.Printf("%c %c\n", goask.ICON_WHITE_RIGHT, result)

	default:
		panic("I don't know that type")
	}

	r.Value = value
	//fmt.Println("You entered", value)
	return r.Value
}
