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

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type Questionaire struct {
	questions []*SmartQuestion
	lastId    *atomic.Uint32
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewQuestionaire() *Questionaire {
	var id atomic.Uint32
	id.Store(0)
	return &Questionaire{
		questions: make([]*SmartQuestion, 0),
		lastId:    &id,
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// add a question and when the answer is obtained, proceed with the
// next question (AskAndContinue mode)
func (qm *Questionaire) AddSequential(q ICurious) uint32 {
	newId := qm.lastId.Add(1)
	sq := NewSmartQuestion(AskAndContinue, q, nil)
	qm.questions = append(qm.questions, sq)
	return newId
}

// add a question and when the answer is obtained, terminate the
// questionaire (AskAndTerminate mode)
func (qm *Questionaire) AddTerminal(q ICurious) uint32 {
	newId := qm.lastId.Add(1)
	sq := NewSmartQuestion(AskAndTerminate, q, nil)
	qm.questions = append(qm.questions, sq)

	return newId
}

// add a multiple choice question and use the callback to determine which would
// be the next question depending on the answer.
func (qm *Questionaire) AddConditionalChoices(q *QuestionWithChoice, callback AnswerCallback) uint32 {
	newId := qm.lastId.Add(1)
	sq := NewSmartQuestion(AskAndDecide, q, callback)
	qm.questions = append(qm.questions, sq)

	return newId
}

// add a smart question which already has the callback information.
func (qm *Questionaire) AddConditionalSmart(q *SmartQuestion) uint32 {
	newId := qm.lastId.Add(1)
	sq := NewSmartQuestion(AskAndDecide, q, q.Callback)
	qm.questions = append(qm.questions, sq)

	return newId
}

// begin the questionaire and terminate when an error occurs or when the
// last question is asked.
func (qm *Questionaire) StartQuestionaire() {
	terminate := false
	nextIdx := 0

	for !terminate {
		// get the current question
		question := qm.questions[nextIdx]
		// ask interactively and get the answer
		question.Ask()
		// decide what to do next
		switch question.Mode {
		// proceed sequentially with the next in the list
		case AskAndContinue:
			nextIdx++

		case AskAndDecide:
			nextQ := question.Next()
			nextIdx = int(nextQ.GetId())

		// terminate the questionaire
		case AskAndTerminate:
			terminate = true
		}
	}
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

/* ----------------------------------------------------------------
 *					M A I N    |     D E M O
 *-----------------------------------------------------------------*/
