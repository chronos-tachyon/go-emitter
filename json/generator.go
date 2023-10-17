package json

import (
	"fmt"
	"math/big"
	"os"

	"github.com/chronos-tachyon/go-emitter"
	"github.com/chronos-tachyon/go-emitter/appenders"
	"github.com/chronos-tachyon/go-emitter/states"
)

type Appender = appenders.Appender

type Generator struct {
	json JSON
	sm   states.Machine
}

func (g *Generator) Reset() {
	g.sm.Reset()
}

func (g *Generator) Factory() emitter.GeneratorFactory {
	return g.json
}

func (g *Generator) Begin() []Appender {
	g.trace("Begin")
	g.sm.ExpectRoot()
	return nil
}

func (g *Generator) End() []Appender {
	g.trace("End")
	g.sm.ExpectEnd()
	g.sm.State = ^states.State(0)

	var b appenders.Builder
	g.lineFeed(&b)
	return b.Build()
}

func (g *Generator) StartObject() []Appender {
	g.trace("StartObject#1")
	g.sm.ExpectValue()
	g.sm.Push(states.ObjectFirstKey)

	var b appenders.Builder
	b.AddByte('{')
	g.trace("StartObject#2")
	return b.Build()
}

func (g *Generator) EndObject() []Appender {
	g.trace("EndObject#1")
	g.sm.ExpectKey()
	needIndent := g.sm.State.In(states.ObjectNextKey)
	g.sm.Pop()

	var b appenders.Builder
	if needIndent {
		g.indent(&b)
	}
	b.AddByte('}')
	g.sm.Next()
	g.trace("EndObject#2")
	return b.Build()
}

func (g *Generator) StartArray() []Appender {
	g.trace("StartArray#1")
	g.sm.ExpectValue()
	g.sm.Push(states.ArrayFirstValue)

	var b appenders.Builder
	b.AddByte('[')
	g.trace("StartArray#2")
	return b.Build()
}

func (g *Generator) EndArray() []Appender {
	g.trace("EndArray#1")
	g.sm.ExpectArray()
	needIndent := g.sm.State.In(states.ArrayNextValue)
	g.sm.Pop()

	var b appenders.Builder
	if needIndent {
		g.indent(&b)
	}
	b.AddByte(']')
	g.sm.Next()
	g.trace("EndArray#2")
	return b.Build()
}

func (g *Generator) Key(key string) []Appender {
	g.trace("Key#1")
	g.sm.ExpectKey()

	var b appenders.Builder
	if g.sm.State.In(states.ObjectFirstKey) {
		g.indent(&b)
	}
	if g.sm.State.In(states.ObjectNextKey) {
		b.AddByte(',')
		g.indentOrSpace(&b)
	}
	b.Add(StringAppender{Value: key, EscapeHTML: g.json.EscapeHTML})
	b.AddByte(':')
	g.space(&b)
	g.sm.Next()
	g.trace("Key#2")
	return b.Build()
}

func (g *Generator) NullValue() []Appender {
	return g.literal(`null`)
}

func (g *Generator) BoolValue(value bool) []Appender {
	if value {
		return g.literal(`true`)
	}
	return g.literal(`false`)
}

func (g *Generator) StringValue(value string) []Appender {
	return g.value(StringAppender{Value: value, EscapeHTML: g.json.EscapeHTML})
}

func (g *Generator) BytesValue(value []byte) []Appender {
	return g.value(BytesAppender{Value: value})
}

func (g *Generator) ByteValue(value byte) []Appender {
	return g.RuneValue(rune(value))
}

func (g *Generator) RuneValue(value rune) []Appender {
	return g.StringValue(string(value))
}

func (g *Generator) IntValue(value int64) []Appender {
	return g.value(appenders.IntText(value))
}

func (g *Generator) UintValue(value uint64) []Appender {
	return g.value(appenders.UintText(value))
}

func (g *Generator) NaNValue() []Appender {
	return g.literal(`"NaN"`)
}

func (g *Generator) InfValue(isNeg bool) []Appender {
	if isNeg {
		return g.literal(`"-Inf"`)
	}
	return g.literal(`"+Inf"`)
}

func (g *Generator) FloatValue(value float64) []Appender {
	return g.value(appenders.FloatText(value))
}

func (g *Generator) BigIntValue(value *big.Int) []Appender {
	if value == nil {
		return g.NullValue()
	}
	return g.value(appenders.BigIntText{Pointer: value})
}

func (g *Generator) BigFloatValue(value *big.Float) []Appender {
	if value == nil {
		return g.NullValue()
	}
	return g.value(appenders.BigFloatText{Pointer: value})
}

func (g *Generator) value(a Appender) []Appender {
	g.trace("value#1")
	g.sm.ExpectValue()

	var b appenders.Builder
	if g.sm.State.In(states.ArrayFirstValue) {
		g.indent(&b)
	}
	if g.sm.State.In(states.ArrayNextValue) {
		b.AddByte(',')
		g.indentOrSpace(&b)
	}
	b.Add(a)
	g.sm.Next()
	g.trace("value#2")
	return b.Build()
}

func (g *Generator) literal(str string) []Appender {
	return g.value(appenders.LiteralString(str))
}

func (g *Generator) indent(b *appenders.Builder) {
	g.json.Format.indent(b, g.json.IndentWithTabs, g.json.IndentSize, g.sm.Depth())
}

func (g *Generator) indentOrSpace(b *appenders.Builder) {
	g.json.Format.indentOrSpace(b, g.json.IndentWithTabs, g.json.IndentSize, g.sm.Depth())
}

func (g *Generator) space(b *appenders.Builder) {
	g.json.Format.space(b)
}

func (g *Generator) lineFeed(b *appenders.Builder) {
	g.json.Format.lineFeed(b)
}

func (g *Generator) trace(call string) {
	if g.json.TraceEnabled {
		fmt.Fprintf(os.Stderr, "%s: json=%v stack=%v state=%v\n", call, g.json, g.sm.Stack, g.sm.State)
	}
}

var _ emitter.Generator = (*Generator)(nil)
