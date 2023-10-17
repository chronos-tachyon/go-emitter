package values

import (
	"math/big"

	"github.com/chronos-tachyon/go-emitter"
)

type Func func(*emitter.Emitter)

func (fn Func) EmitTo(e *emitter.Emitter) {
	fn(e)
}

var _ emitter.Value = Func(nil)

type Null struct{}

func (Null) EmitTo(e *emitter.Emitter) {
	e.EmitNull()
}

var _ emitter.Value = Null{}

type Bool bool

func (v Bool) EmitTo(e *emitter.Emitter) {
	e.EmitBool(bool(v))
}

var _ emitter.Value = Bool(false)

type Int int64

func (v Int) EmitTo(e *emitter.Emitter) {
	e.EmitInt64(int64(v))
}

var _ emitter.Value = Int(0)

type Uint uint64

func (v Uint) EmitTo(e *emitter.Emitter) {
	e.EmitUint64(uint64(v))
}

var _ emitter.Value = Uint(0)

type BigIntValue struct{ Pointer *big.Int }

func (v BigIntValue) EmitTo(e *emitter.Emitter) {
	e.EmitBigInt(v.Pointer)
}

var _ emitter.Value = BigIntValue{}

type Float float64

func (v Float) EmitTo(e *emitter.Emitter) {
	e.EmitFloat64(float64(v))
}

var _ emitter.Value = Float(0)

type BigFloatValue struct{ Pointer *big.Float }

func (v BigFloatValue) EmitTo(e *emitter.Emitter) {
	e.EmitBigFloat(v.Pointer)
}

var _ emitter.Value = BigFloatValue{}

type String string

func (v String) EmitTo(e *emitter.Emitter) {
	e.EmitString(string(v))
}

var _ emitter.Value = String("")

type Bytes []byte

func (v Bytes) EmitTo(e *emitter.Emitter) {
	e.EmitBytes([]byte(v))
}

var _ emitter.Value = Bytes(nil)

type Byte byte

func (v Byte) EmitTo(e *emitter.Emitter) {
	e.EmitByte(byte(v))
}

var _ emitter.Value = Byte('a')

type Rune rune

func (v Rune) EmitTo(e *emitter.Emitter) {
	e.EmitRune(rune(v))
}

var _ emitter.Value = Rune('a')

type Array []emitter.Value

func (v Array) EmitTo(e *emitter.Emitter) {
	e.StartArray()
	for _, item := range v {
		item.EmitTo(e)
	}
	e.EndArray()
}

var _ emitter.Value = Array(nil)

type ObjectField struct {
	Key   string
	Value emitter.Value
}

type Object []ObjectField

func (v Object) EmitTo(e *emitter.Emitter) {
	e.StartObject()
	for _, item := range v {
		e.EmitKey(item.Key)
		item.Value.EmitTo(e)
	}
	e.EndObject()
}

var _ emitter.Value = Object(nil)
