package emitter

import (
	"math/big"
)

type Value interface {
	EmitTo(*Emitter)
}

type FuncValue func(*Emitter)

func (fn FuncValue) EmitTo(e *Emitter) {
	fn(e)
}

var _ Value = FuncValue(nil)

type NullValue struct{}

func (NullValue) EmitTo(e *Emitter) {
	e.EmitNull()
}

var _ Value = NullValue{}

type BoolValue bool

func (v BoolValue) EmitTo(e *Emitter) {
	e.EmitBool(bool(v))
}

var _ Value = BoolValue(false)

type IntValue int64

func (v IntValue) EmitTo(e *Emitter) {
	e.EmitInt64(int64(v))
}

var _ Value = IntValue(0)

type UintValue uint64

func (v UintValue) EmitTo(e *Emitter) {
	e.EmitUint64(uint64(v))
}

var _ Value = UintValue(0)

type BigIntValue struct{ Pointer *big.Int }

func (v BigIntValue) EmitTo(e *Emitter) {
	e.EmitBigInt(v.Pointer)
}

var _ Value = BigIntValue{}

type FloatValue float64

func (v FloatValue) EmitTo(e *Emitter) {
	e.EmitFloat64(float64(v))
}

var _ Value = FloatValue(0)

type BigFloatValue struct{ Pointer *big.Float }

func (v BigFloatValue) EmitTo(e *Emitter) {
	e.EmitBigFloat(v.Pointer)
}

var _ Value = BigFloatValue{}

type StringValue string

func (v StringValue) EmitTo(e *Emitter) {
	e.EmitString(string(v))
}

var _ Value = StringValue("")

type BytesValue []byte

func (v BytesValue) EmitTo(e *Emitter) {
	e.EmitBytes([]byte(v))
}

var _ Value = BytesValue(nil)

type ByteValue byte

func (v ByteValue) EmitTo(e *Emitter) {
	e.EmitByte(byte(v))
}

var _ Value = ByteValue('a')

type RuneValue rune

func (v RuneValue) EmitTo(e *Emitter) {
	e.EmitRune(rune(v))
}

var _ Value = RuneValue('a')

type ArrayValue []Value

func (v ArrayValue) EmitTo(e *Emitter) {
	e.StartArray()
	for _, item := range v {
		item.EmitTo(e)
	}
	e.EndArray()
}

var _ Value = ArrayValue(nil)

type ObjectField struct {
	Key   string
	Value Value
}

type ObjectValue []ObjectField

func (v ObjectValue) EmitTo(e *Emitter) {
	e.StartObject()
	for _, item := range v {
		e.EmitKey(item.Key)
		item.Value.EmitTo(e)
	}
	e.EndObject()
}

var _ Value = ObjectValue(nil)
