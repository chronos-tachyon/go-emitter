package emitter

import (
	"math/big"
)

type GeneratorFactory interface {
	NewGenerator() Generator
}

type Generator interface {
	Reset()
	Factory() GeneratorFactory

	Begin() []Appender
	End() []Appender

	StartObject() []Appender
	EndObject() []Appender

	StartArray() []Appender
	EndArray() []Appender

	Key(key string) []Appender

	NullValue() []Appender

	BoolValue(value bool) []Appender

	IntValue(value int64) []Appender
	UintValue(value uint64) []Appender
	BigIntValue(value *big.Int) []Appender

	NaNValue() []Appender
	InfValue(isNeg bool) []Appender
	FloatValue(value float64) []Appender
	BigFloatValue(value *big.Float) []Appender

	StringValue(value string) []Appender
	BytesValue(value []byte) []Appender
	ByteValue(value byte) []Appender
	RuneValue(value rune) []Appender
}
