/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package ask

import "sync/atomic"

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	AskAndContinue AnswerType = iota
	AskAndDecide
	AskAndTerminate
)

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ ICuriouslySmart = (*SmartQuestion)(nil)
var smartIdGenerator atomic.Uint32

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type AnswerType uint8

type AnswerCallback func(id uint32) *SmartQuestion

type SmartQuestion struct {
	Id       uint32
	Mode     AnswerType
	Question ICurious
	Callback AnswerCallback
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewSmartQuestion(mode AnswerType, question ICurious, callback AnswerCallback) *SmartQuestion {
	id := smartIdGenerator.Add(1)
	if mode == AskAndTerminate || mode == AskAndContinue {
		callback = nil
	} else if callback == nil {
		mode = AskAndTerminate
	}

	return &SmartQuestion{
		Id:       id,
		Mode:     mode,
		Question: question,
		Callback: callback,
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

func (q *SmartQuestion) GetId() uint32 {
	return q.Id
}

func (q *SmartQuestion) Type() AnswerType {
	return q.Mode
}

func (q *SmartQuestion) Next() ICuriouslySmart {
	if q.Callback != nil {
		return q.Callback(uint32(q.Question.AsInt()))
	}
	return nil
}

func (q *SmartQuestion) Answer() any {
	return q.Question.Answer()
}

func (q *SmartQuestion) Ask() ICurious {
	return q.Question.Ask()
}

func (q *SmartQuestion) AsInt() int {
	if v, ok := q.Answer().(int); ok {
		return v
	}
	println("The answer is not an integer")
	return -1
}

func (q *SmartQuestion) AsRune() rune {
	if v, ok := q.Answer().(rune); ok {
		return v
	}
	println("The answer is not a rune")
	return rune(0)
}

func (q *SmartQuestion) AsString() string {
	if v, ok := q.Answer().(string); ok {
		return v
	}
	println("The answer is not a string")
	return ""
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						T e s t s
 *-----------------------------------------------------------------*/
