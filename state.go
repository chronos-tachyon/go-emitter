package emitter

import (
	"fmt"
)

type State byte

const (
	RootState State = iota
	ObjectFirstKeyState
	ObjectNextKeyState
	ObjectFirstValueState
	ObjectNextValueState
	ArrayFirstValueState
	ArrayNextValueState
	EndState
)

const emitterStateSize = 8

var emitterStateGoNames = [emitterStateSize]string{
	"emitter.RootState",
	"emitter.ObjectFirstKeyState",
	"emitter.ObjectNextKeyState",
	"emitter.ObjectFirstValueState",
	"emitter.ObjectNextValueState",
	"emitter.ArrayFirstValueState",
	"emitter.ArrayNextValueState",
	"emitter.EndState",
}

var emitterStateNames = [emitterStateSize]string{
	"root",
	"objectFirstKey",
	"objectNextKey",
	"objectFirstValue",
	"objectNextValue",
	"arrayFirstValue",
	"arrayNextValue",
	"end",
}

func (state State) IsValid() bool {
	return state < emitterStateSize
}

func (state State) GoString() string {
	if state.IsValid() {
		return emitterStateGoNames[state]
	}
	return fmt.Sprintf("emitter.State(%d)", uint(state))
}

func (state State) String() string {
	if state.IsValid() {
		return emitterStateNames[state]
	}
	return fmt.Sprintf("%%!ERR[invalid emitter.State %d]", uint(state))
}

func (state State) Next() State {
	switch state {
	case RootState:
		return EndState

	case ObjectFirstKeyState:
		return ObjectFirstValueState
	case ObjectNextKeyState:
		return ObjectNextValueState

	case ObjectFirstValueState:
		fallthrough
	case ObjectNextValueState:
		return ObjectNextKeyState

	case ArrayFirstValueState:
		fallthrough
	case ArrayNextValueState:
		return ArrayNextValueState

	default:
		panic(state.Unexpected())
	}
}

func (state State) In(oneOf ...State) bool {
	for _, x := range oneOf {
		if state == x {
			return true
		}
	}
	return false
}

func (state State) Expect(oneOf ...State) {
	if !state.In(oneOf...) {
		panic(state.Unexpected())
	}
}

func (state State) Unexpected(oneOf ...State) error {
	if len(oneOf) <= 0 {
		return fmt.Errorf("unexpected state %v", state)
	}
	return fmt.Errorf("unexpected state %v; expected one of %v", state, oneOf)
}

var (
	_ fmt.GoStringer = State(0)
	_ fmt.Stringer   = State(0)
)
