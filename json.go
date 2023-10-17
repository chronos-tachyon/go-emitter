package emitter

import (
	"fmt"
	"math/big"
	"os"
)

var TracingEnabled bool = false

type JSON struct {
	Format         JSONFormat
	IndentSize     uint
	IndentWithTabs bool
	EscapeHTML     bool
}

func (json JSON) New() Generator {
	g := &jsonGenerator{json: json}
	g.Reset()
	return g
}

var _ GeneratorFactory = JSON{}

type jsonGenerator struct {
	json JSON
	sm   StateMachine
}

func (g *jsonGenerator) Reset() {
	g.sm.Reset()
}

func (g *jsonGenerator) Factory() GeneratorFactory {
	return g.json
}

func (g *jsonGenerator) End() []Appender {
	g.trace("End")
	g.sm.ExpectEnd()
	g.sm.State = ^State(0)

	var b AppenderBuilder
	g.lineFeed(&b)
	return b.Build()
}

func (g *jsonGenerator) StartObject() []Appender {
	g.trace("StartObject#1")
	g.sm.ExpectValue()
	g.sm.Push(ObjectFirstKeyState)

	var b AppenderBuilder
	b.AddByte('{')
	g.trace("StartObject#2")
	return b.Build()
}

func (g *jsonGenerator) EndObject() []Appender {
	g.trace("EndObject#1")
	g.sm.ExpectKey()
	needIndent := g.sm.State.In(ObjectNextKeyState)
	g.sm.Pop()

	var b AppenderBuilder
	if needIndent {
		g.indent(&b)
	}
	b.AddByte('}')
	g.sm.Next()
	g.trace("EndObject#2")
	return b.Build()
}

func (g *jsonGenerator) StartArray() []Appender {
	g.trace("StartArray#1")
	g.sm.ExpectValue()
	g.sm.Push(ArrayFirstValueState)

	var b AppenderBuilder
	b.AddByte('[')
	g.trace("StartArray#2")
	return b.Build()
}

func (g *jsonGenerator) EndArray() []Appender {
	g.trace("EndArray#1")
	g.sm.ExpectArray()
	needIndent := g.sm.State.In(ArrayNextValueState)
	g.sm.Pop()

	var b AppenderBuilder
	if needIndent {
		g.indent(&b)
	}
	b.AddByte(']')
	g.sm.Next()
	g.trace("EndArray#2")
	return b.Build()
}

func (g *jsonGenerator) Key(key string) []Appender {
	g.trace("Key#1")
	g.sm.ExpectKey()

	var b AppenderBuilder
	if g.sm.State.In(ObjectFirstKeyState) {
		g.indent(&b)
	}
	if g.sm.State.In(ObjectNextKeyState) {
		b.AddByte(',')
		g.indentOrSpace(&b)
	}
	b.Add(&jsonAppenderString{value: key, escapeHTML: g.json.EscapeHTML})
	b.AddByte(':')
	g.space(&b)
	g.sm.Next()
	g.trace("Key#2")
	return b.Build()
}

func (g *jsonGenerator) NullValue() []Appender {
	return g.literal(`null`)
}

func (g *jsonGenerator) BoolValue(value bool) []Appender {
	if value {
		return g.literal(`true`)
	}
	return g.literal(`false`)
}

func (g *jsonGenerator) StringValue(value string) []Appender {
	return g.value(&jsonAppenderString{value: value, escapeHTML: g.json.EscapeHTML})
}

func (g *jsonGenerator) BytesValue(value []byte) []Appender {
	return g.value(&jsonAppenderBytes{value: value})
}

func (g *jsonGenerator) ByteValue(value byte) []Appender {
	return g.RuneValue(rune(value))
}

func (g *jsonGenerator) RuneValue(value rune) []Appender {
	return g.StringValue(string(value))
}

func (g *jsonGenerator) IntValue(value int64) []Appender {
	return g.value(IntTextAppender(value))
}

func (g *jsonGenerator) UintValue(value uint64) []Appender {
	return g.value(UintTextAppender(value))
}

func (g *jsonGenerator) NaNValue() []Appender {
	return g.literal(`"NaN"`)
}

func (g *jsonGenerator) InfValue(isNeg bool) []Appender {
	if isNeg {
		return g.literal(`"-Inf"`)
	}
	return g.literal(`"+Inf"`)
}

func (g *jsonGenerator) FloatValue(value float64) []Appender {
	return g.value(FloatTextAppender(value))
}

func (g *jsonGenerator) BigIntValue(value *big.Int) []Appender {
	if value == nil {
		return g.NullValue()
	}
	return g.value(BigIntTextAppender{Pointer: value})
}

func (g *jsonGenerator) BigFloatValue(value *big.Float) []Appender {
	if value == nil {
		return g.NullValue()
	}
	return g.value(BigFloatTextAppender{Pointer: value})
}

func (g *jsonGenerator) value(a Appender) []Appender {
	g.trace("value#1")
	g.sm.ExpectValue()

	var b AppenderBuilder
	if g.sm.State.In(ArrayFirstValueState) {
		g.indent(&b)
	}
	if g.sm.State.In(ArrayNextValueState) {
		b.AddByte(',')
		g.indentOrSpace(&b)
	}
	b.Add(a)
	g.sm.Next()
	g.trace("value#2")
	return b.Build()
}

func (g *jsonGenerator) literal(str string) []Appender {
	return g.value(LiteralStringAppender(str))
}

func (g *jsonGenerator) indent(b *AppenderBuilder) {
	g.json.Format.indent(b, g.json.IndentWithTabs, g.json.IndentSize, g.sm.Depth())
}

func (g *jsonGenerator) indentOrSpace(b *AppenderBuilder) {
	g.json.Format.indentOrSpace(b, g.json.IndentWithTabs, g.json.IndentSize, g.sm.Depth())
}

func (g *jsonGenerator) space(b *AppenderBuilder) {
	g.json.Format.space(b)
}

func (g *jsonGenerator) lineFeed(b *AppenderBuilder) {
	g.json.Format.lineFeed(b)
}

func (g *jsonGenerator) trace(call string) {
	if TracingEnabled {
		fmt.Fprintf(os.Stderr, "%s: json=%v stack=%v state=%v\n", call, g.json, g.sm.Stack, g.sm.State)
	}
}

var _ Generator = (*jsonGenerator)(nil)
