package states

import (
	"fmt"
)

type State byte

const (
	Root State = iota
	ObjectFirstKey
	ObjectNextKey
	ObjectFirstValue
	ObjectNextValue
	ArrayFirstValue
	ArrayNextValue
	End
)

const stateSize = 8

var stateGoNames = [stateSize]string{
	"states.Root",
	"states.ObjectFirstKey",
	"states.ObjectNextKey",
	"states.ObjectFirstValue",
	"states.ObjectNextValue",
	"states.ArrayFirstValue",
	"states.ArrayNextValue",
	"states.End",
}

var stateNames = [stateSize]string{
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
	return state < stateSize
}

func (state State) GoString() string {
	if state.IsValid() {
		return stateGoNames[state]
	}
	return fmt.Sprintf("states.State(%d)", uint(state))
}

func (state State) String() string {
	if state.IsValid() {
		return stateNames[state]
	}
	return fmt.Sprintf("%%!ERR[invalid states.State %d]", uint(state))
}

func (state State) Next() State {
	switch state {
	case Root:
		return End

	case ObjectFirstKey:
		return ObjectFirstValue
	case ObjectNextKey:
		return ObjectNextValue

	case ObjectFirstValue:
		fallthrough
	case ObjectNextValue:
		return ObjectNextKey

	case ArrayFirstValue:
		fallthrough
	case ArrayNextValue:
		return ArrayNextValue

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
