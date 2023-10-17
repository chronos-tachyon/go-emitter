package emitter

import (
	"fmt"
	"io"
	"math"
	"math/big"
	"reflect"
)

const (
	blockSize  = 4096
	bufferSize = blockSize * 2
)

type Emitter struct {
	w       io.Writer
	g       Generator
	out     []byte
	err     error
	n       int64
	scratch [bufferSize]byte
}

func New(w io.Writer, g Generator) *Emitter {
	e := &Emitter{}
	e.Reset(w, g)
	return e
}

func (e *Emitter) Reset(w io.Writer, g Generator) {
	if w == nil {
		panic(fmt.Errorf("io.Writer is nil"))
	}

	if g == nil {
		panic(fmt.Errorf("emitter.Generator is nil"))
	}

	g.Reset()
	*e = Emitter{w: w, g: g}
	e.out = e.scratch[:0]
	e.apply(g.Begin())
}

func (e *Emitter) Writer() io.Writer {
	return e.w
}

func (e *Emitter) Generator() Generator {
	return e.g
}

func (e *Emitter) BytesWritten() int64 {
	return e.n
}

func (e *Emitter) StartObject() {
	e.apply(e.g.StartObject())
}

func (e *Emitter) EndObject() {
	e.apply(e.g.EndObject())
}

func (e *Emitter) StartArray() {
	e.apply(e.g.StartArray())
}

func (e *Emitter) EndArray() {
	e.apply(e.g.EndArray())
}

func (e *Emitter) EmitKey(key string) {
	e.apply(e.g.Key(key))
}

func (e *Emitter) EmitValue(value Value) {
	value.EmitTo(e)
}

func (e *Emitter) EmitNull() {
	e.apply(e.g.NullValue())
}

func (e *Emitter) EmitBool(value bool) {
	e.apply(e.g.BoolValue(value))
}

func (e *Emitter) EmitInt(value int) {
	e.EmitInt64(int64(value))
}

func (e *Emitter) EmitInt8(value int8) {
	e.EmitInt64(int64(value))
}

func (e *Emitter) EmitInt16(value int16) {
	e.EmitInt64(int64(value))
}

func (e *Emitter) EmitInt32(value int32) {
	e.EmitInt64(int64(value))
}

func (e *Emitter) EmitInt64(value int64) {
	e.apply(e.g.IntValue(value))
}

func (e *Emitter) EmitUint(value uint) {
	e.EmitUint64(uint64(value))
}

func (e *Emitter) EmitUint8(value uint8) {
	e.EmitUint64(uint64(value))
}

func (e *Emitter) EmitUint16(value uint16) {
	e.EmitUint64(uint64(value))
}

func (e *Emitter) EmitUint32(value uint32) {
	e.EmitUint64(uint64(value))
}

func (e *Emitter) EmitUint64(value uint64) {
	e.apply(e.g.UintValue(value))
}

func (e *Emitter) EmitBigInt(value *big.Int) {
	e.apply(e.g.BigIntValue(value))
}

func (e *Emitter) EmitFloat32(value float32) {
	e.EmitFloat64(float64(value))
}

func (e *Emitter) EmitFloat64(value float64) {
	switch {
	case math.IsNaN(value):
		e.apply(e.g.NaNValue())
	case math.IsInf(value, 0):
		e.apply(e.g.InfValue(value < 0))
	default:
		e.apply(e.g.FloatValue(value))
	}
}

func (e *Emitter) EmitBigFloat(value *big.Float) {
	e.apply(e.g.BigFloatValue(value))
}

func (e *Emitter) EmitString(value string) {
	e.apply(e.g.StringValue(value))
}

func (e *Emitter) EmitBytes(value []byte) {
	e.apply(e.g.BytesValue(value))
}

func (e *Emitter) EmitByte(value byte) {
	e.apply(e.g.ByteValue(value))
}

func (e *Emitter) EmitRune(value rune) {
	e.apply(e.g.RuneValue(value))
}

func (e *Emitter) Emit(value any) {
	switch x := value.(type) {
	case nil:
		e.EmitNull()
	case Value:
		e.EmitValue(x)
	case bool:
		e.EmitBool(x)
	case int64:
		e.EmitInt64(x)
	case int32:
		e.EmitInt32(x)
	case int16:
		e.EmitInt16(x)
	case int8:
		e.EmitInt8(x)
	case int:
		e.EmitInt(x)
	case uint64:
		e.EmitUint64(x)
	case uint32:
		e.EmitUint32(x)
	case uint16:
		e.EmitUint16(x)
	case uint8:
		e.EmitUint8(x)
	case uint:
		e.EmitUint(x)
	case *big.Int:
		e.EmitBigInt(x)
	case float64:
		e.EmitFloat64(x)
	case float32:
		e.EmitFloat32(x)
	case *big.Float:
		e.EmitBigFloat(x)
	case string:
		e.EmitString(x)
	case []byte:
		e.EmitBytes(x)
	case reflect.Value:
		e.EmitReflected(x)
	default:
		e.EmitReflected(reflect.ValueOf(value))
	}
}

func (e *Emitter) EmitReflected(value reflect.Value) {
	panic(fmt.Errorf("EmitReflected not implemented for %v", value.Type()))
}

func (e *Emitter) Flush() error {
	if e.err == nil && len(e.out) > 0 {
		e.flush()
	}
	return nil
}

func (e *Emitter) Close() error {
	e.apply(e.g.End())
	if e.err == nil {
		e.flush()
	}
	return e.err
}

func (e *Emitter) apply(list []Appender) {
	if e.err != nil {
		return
	}

	for _, a := range list {
		e.out = a.Append(e.out)
	}

	if len(e.out) < blockSize {
		return
	}

	e.flush()
}

func (e *Emitter) flush() {
	var n int
	n, e.err = e.w.Write(e.out)
	e.n += int64(n)
	e.out = e.scratch[:0]
}
