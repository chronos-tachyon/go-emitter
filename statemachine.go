package emitter

import (
	"fmt"
)

const kStackSize = 16

type StateMachine struct {
	Stack []State
	State State

	scratch [kStackSize]State
}

func (sm *StateMachine) Reset() {
	*sm = StateMachine{}
	sm.Stack = sm.scratch[:0]
}

func (sm *StateMachine) Depth() uint {
	return uint(len(sm.Stack))
}

func (sm *StateMachine) Push(next State) {
	sm.Stack = append(sm.Stack, sm.State)
	sm.State = next
}

func (sm *StateMachine) Pop() {
	n := uint(len(sm.Stack))
	if n <= 0 {
		panic(fmt.Errorf("stack is empty"))
	}

	n--
	sm.State = sm.Stack[n]
	sm.Stack = sm.Stack[:n]
}

func (sm *StateMachine) Next() {
	sm.State = sm.State.Next()
}

func (sm *StateMachine) Expect(oneOf ...State) {
	sm.State.Expect(oneOf...)
}

func (sm *StateMachine) ExpectKey() {
	sm.Expect(ObjectFirstKeyState, ObjectNextKeyState)
}

func (sm *StateMachine) ExpectValue() {
	sm.Expect(RootState, ObjectFirstValueState, ObjectNextValueState, ArrayFirstValueState, ArrayNextValueState)
}

func (sm *StateMachine) ExpectArray() {
	sm.Expect(ArrayFirstValueState, ArrayNextValueState)
}

func (sm *StateMachine) ExpectEnd() {
	sm.Expect(EndState)
}
