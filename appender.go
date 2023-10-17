package emitter

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"strconv"
)

type Appender interface {
	Append(out []byte) []byte
}

type AppenderBuilder struct {
	list []Appender
	out  []byte
}

func (b *AppenderBuilder) flush() {
	if len(b.out) <= 0 {
		return
	}
	str := string(b.out)
	b.out = nil
	b.list = append(b.list, LiteralStringAppender(str))
}

func (b *AppenderBuilder) Add(a Appender) {
	switch x := a.(type) {
	case nil:
		// pass
	case LiteralBytesAppender:
		b.AddBytes([]byte(x))
	case LiteralStringAppender:
		b.AddString(string(x))
	default:
		b.flush()
		b.list = append(b.list, a)
	}
}

func (b *AppenderBuilder) Indent(tab bool, size uint, count uint) {
	if count <= 0 {
		return
	}

	unit, defaultSize := byte(' '), uint(2)
	if tab {
		unit, defaultSize = '\t', 1
	}
	if size <= 0 {
		size = defaultSize
	}

	numUnits := count * size
	for i := uint(0); i < numUnits; i++ {
		b.out = append(b.out, unit)
	}
}

func (b *AppenderBuilder) AddBytes(buf []byte) {
	b.out = append(b.out, buf...)
}

func (b *AppenderBuilder) AddString(str string) {
	b.out = append(b.out, str...)
}

func (b *AppenderBuilder) AddByte(ch byte) {
	b.out = append(b.out, ch)
}

func (b *AppenderBuilder) AddRune(ch rune) {
	str := string(ch)
	b.out = append(b.out, str...)
}

func (b *AppenderBuilder) Build() []Appender {
	b.flush()
	list := b.list
	b.list = nil
	return list
}

type LiteralBytesAppender []byte

func (a LiteralBytesAppender) Append(out []byte) []byte {
	slice := []byte(a)
	return append(out, slice...)
}

var _ Appender = LiteralBytesAppender(nil)

type LiteralStringAppender string

func (a LiteralStringAppender) Append(out []byte) []byte {
	str := string(a)
	return append(out, str...)
}

var _ Appender = LiteralStringAppender("")

type IntTextAppender int64

func (a IntTextAppender) Append(out []byte) []byte {
	i64 := int64(a)
	return strconv.AppendInt(out, i64, 10)
}

var _ Appender = IntTextAppender(0)

type UintTextAppender uint64

func (a UintTextAppender) Append(out []byte) []byte {
	u64 := uint64(a)
	return strconv.AppendUint(out, u64, 10)
}

var _ Appender = UintTextAppender(0)

type BigIntTextAppender struct {
	Pointer *big.Int
}

func (a BigIntTextAppender) Append(out []byte) []byte {
	return a.Pointer.Append(out, 10)
}

var _ Appender = BigIntTextAppender{}

type FloatTextAppender float64

func (a FloatTextAppender) Append(out []byte) []byte {
	f64 := float64(a)
	return strconv.AppendFloat(out, f64, 'g', -1, 64)
}

var _ Appender = FloatTextAppender(0)

type BigFloatTextAppender struct {
	Pointer *big.Float
}

func (a BigFloatTextAppender) Append(out []byte) []byte {
	return a.Pointer.Append(out, 'g', -1)
}

var _ Appender = BigFloatTextAppender{}

type jsonAppenderString struct {
	value      string
	escapeHTML bool
}

func (a *jsonAppenderString) Append(out []byte) []byte {
	out = append(out, '"')
	for _, ch := range a.value {
		if esc, found := stringEscapes[ch]; found {
			out = append(out, esc...)
			continue
		}

		if a.escapeHTML {
			if esc, found := stringEscapesHTML[ch]; found {
				out = append(out, esc...)
				continue
			}
		}

		if ch >= 0xd800 && ch < 0xe000 {
			out = fmt.Appendf(out, "\\u%04x", ch)
			continue
		}

		str := string(ch)
		out = append(out, str...)
	}
	return append(out, '"')
}

var _ Appender = (*jsonAppenderString)(nil)

type jsonAppenderBytes struct {
	value []byte
}

func (a *jsonAppenderBytes) Append(out []byte) []byte {
	out = append(out, '"')
	n := uint(len(a.value))
	for i := uint(0); i < n; i += 3 {
		j := min(i+3, n)
		var tmp [4]byte
		base64.StdEncoding.Encode(tmp[:], a.value[i:j])
		out = append(out, tmp[:]...)
	}
	return append(out, '"')
}

var _ Appender = (*jsonAppenderBytes)(nil)

var stringEscapes = map[rune]string{
	'\b': `\b`,
	'\f': `\f`,
	'\n': `\n`,
	'\r': `\r`,
	'\t': `\t`,
	'"':  `\"`,
	'\\': `\\`,
}

var stringEscapesHTML = map[rune]string{
	'&': `\u0026`,
	'<': `\u003c`,
	'>': `\u003e`,
}

func init() {
	for ch := rune(0); ch < 0x20; ch++ {
		if _, found := stringEscapes[ch]; !found {
			stringEscapes[ch] = fmt.Sprintf("\\u%04x", ch)
		}
	}
}
