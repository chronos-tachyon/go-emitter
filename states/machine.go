package states

import (
	"fmt"
)

const stackSize = 16

type Machine struct {
	Stack []State
	State State

	scratch [stackSize]State
}

func (sm *Machine) Reset() {
	*sm = Machine{}
	sm.Stack = sm.scratch[:0]
}

func (sm *Machine) Depth() uint {
	return uint(len(sm.Stack))
}

func (sm *Machine) Push(next State) {
	sm.Stack = append(sm.Stack, sm.State)
	sm.State = next
}

func (sm *Machine) Pop() {
	n := uint(len(sm.Stack))
	if n <= 0 {
		panic(fmt.Errorf("stack is empty"))
	}

	n--
	sm.State = sm.Stack[n]
	sm.Stack = sm.Stack[:n]
}

func (sm *Machine) Next() {
	sm.State = sm.State.Next()
}

func (sm *Machine) Expect(oneOf ...State) {
	sm.State.Expect(oneOf...)
}

func (sm *Machine) ExpectRoot() {
	sm.Expect(Root)
}

func (sm *Machine) ExpectKey() {
	sm.Expect(ObjectFirstKey, ObjectNextKey)
}

func (sm *Machine) ExpectValue() {
	sm.Expect(Root, ObjectFirstValue, ObjectNextValue, ArrayFirstValue, ArrayNextValue)
}

func (sm *Machine) ExpectArray() {
	sm.Expect(ArrayFirstValue, ArrayNextValue)
}

func (sm *Machine) ExpectEnd() {
	sm.Expect(End)
}
