package json

import (
	"bytes"
	"reflect"
	"testing"
	"unicode/utf8"

	"github.com/chronos-tachyon/go-emitter"
	"github.com/chronos-tachyon/go-emitter/values"
)

type Value = emitter.Value

var (
	kNullValue   Value = values.Null{}
	kFalseValue  Value = values.Bool(false)
	kTrueValue   Value = values.Bool(true)
	kIntValue    Value = values.Int(-42)
	kUintValue   Value = values.Uint(42)
	kFloatValue  Value = values.Float(4.2)
	kStringValue Value = values.String("abc")
	kBytesValue  Value = values.Bytes("abc")

	kArrayValue Value = values.Array{
		values.Rune('a'),
		values.Rune('b'),
		values.Rune('c'),
	}

	kObjectValue Value = values.Object{
		{"a", values.Int(1)},
		{"b", values.Int(2)},
		{"c", values.Int(3)},
	}

	kFancyValue Value = values.Object{
		{"@type", values.String("Foo")},
		{"emptyList", values.Array(nil)},
		{"emptyObject", values.Object(nil)},
		{"array", kArrayValue},
		{"object", kObjectValue},
	}
)

func TestJSON(t *testing.T) {
	type testCase struct {
		Name      string
		Input     Value
		Factory   emitter.GeneratorFactory
		Expect    []byte
		ExpectErr error
	}

	compactJSON := JSON{}
	oneLineJSON := JSON{Format: OneLine}
	multiLineJSON := JSON{Format: MultiLine}

	testData := [...]testCase{
		{
			Name:    "Compact/Null",
			Input:   kNullValue,
			Factory: compactJSON,
			Expect:  []byte(`null`),
		},
		{
			Name:    "Compact/False",
			Input:   kFalseValue,
			Factory: compactJSON,
			Expect:  []byte(`false`),
		},
		{
			Name:    "Compact/True",
			Input:   kTrueValue,
			Factory: compactJSON,
			Expect:  []byte(`true`),
		},
		{
			Name:    "Compact/Int",
			Input:   kIntValue,
			Factory: compactJSON,
			Expect:  []byte(`-42`),
		},
		{
			Name:    "Compact/Uint",
			Input:   kUintValue,
			Factory: compactJSON,
			Expect:  []byte(`42`),
		},
		{
			Name:    "Compact/Float",
			Input:   kFloatValue,
			Factory: compactJSON,
			Expect:  []byte(`4.2`),
		},
		{
			Name:    "Compact/String",
			Input:   kStringValue,
			Factory: compactJSON,
			Expect:  []byte(`"abc"`),
		},
		{
			Name:    "Compact/Bytes",
			Input:   kBytesValue,
			Factory: compactJSON,
			Expect:  []byte(`"YWJj"`),
		},
		{
			Name:    "Compact/Array",
			Input:   kArrayValue,
			Factory: compactJSON,
			Expect:  []byte(`["a","b","c"]`),
		},
		{
			Name:    "Compact/Object",
			Input:   kObjectValue,
			Factory: compactJSON,
			Expect:  []byte(`{"a":1,"b":2,"c":3}`),
		},
		{
			Name:    "Compact/Fancy",
			Input:   kFancyValue,
			Factory: compactJSON,
			Expect:  []byte(`{"@type":"Foo","emptyList":[],"emptyObject":{},"array":["a","b","c"],"object":{"a":1,"b":2,"c":3}}`),
		},

		{
			Name:    "OneLine/Null",
			Input:   kNullValue,
			Factory: oneLineJSON,
			Expect:  ParseOneLine(`null`),
		},
		{
			Name:    "OneLine/False",
			Input:   kFalseValue,
			Factory: oneLineJSON,
			Expect:  ParseOneLine(`false`),
		},
		{
			Name:    "OneLine/True",
			Input:   kTrueValue,
			Factory: oneLineJSON,
			Expect:  ParseOneLine(`true`),
		},
		{
			Name:    "OneLine/Int",
			Input:   kIntValue,
			Factory: oneLineJSON,
			Expect:  ParseOneLine(`-42`),
		},
		{
			Name:    "OneLine/Uint",
			Input:   kUintValue,
			Factory: oneLineJSON,
			Expect:  ParseOneLine(`42`),
		},
		{
			Name:    "OneLine/Float",
			Input:   kFloatValue,
			Factory: oneLineJSON,
			Expect:  ParseOneLine(`4.2`),
		},
		{
			Name:    "OneLine/String",
			Input:   kStringValue,
			Factory: oneLineJSON,
			Expect:  ParseOneLine(`"abc"`),
		},
		{
			Name:    "OneLine/Bytes",
			Input:   kBytesValue,
			Factory: oneLineJSON,
			Expect:  ParseOneLine(`"YWJj"`),
		},
		{
			Name:    "OneLine/Array",
			Input:   kArrayValue,
			Factory: oneLineJSON,
			Expect:  ParseOneLine(`["a", "b", "c"]`),
		},
		{
			Name:    "OneLine/Object",
			Input:   kObjectValue,
			Factory: oneLineJSON,
			Expect:  ParseOneLine(`{"a": 1, "b": 2, "c": 3}`),
		},
		{
			Name:    "OneLine/Fancy",
			Input:   kFancyValue,
			Factory: oneLineJSON,
			Expect:  ParseOneLine(`{"@type": "Foo", "emptyList": [], "emptyObject": {}, "array": ["a", "b", "c"], "object": {"a": 1, "b": 2, "c": 3}}`),
		},

		{
			Name:    "MultiLine/Null",
			Input:   kNullValue,
			Factory: multiLineJSON,
			Expect:  ParseOneLine(`null`),
		},
		{
			Name:    "MultiLine/False",
			Input:   kFalseValue,
			Factory: multiLineJSON,
			Expect:  ParseOneLine(`false`),
		},
		{
			Name:    "MultiLine/True",
			Input:   kTrueValue,
			Factory: multiLineJSON,
			Expect:  ParseOneLine(`true`),
		},
		{
			Name:    "MultiLine/Int",
			Input:   kIntValue,
			Factory: multiLineJSON,
			Expect:  ParseOneLine(`-42`),
		},
		{
			Name:    "MultiLine/Uint",
			Input:   kUintValue,
			Factory: multiLineJSON,
			Expect:  ParseOneLine(`42`),
		},
		{
			Name:    "MultiLine/Float",
			Input:   kFloatValue,
			Factory: multiLineJSON,
			Expect:  ParseOneLine(`4.2`),
		},
		{
			Name:    "MultiLine/String",
			Input:   kStringValue,
			Factory: multiLineJSON,
			Expect:  ParseOneLine(`"abc"`),
		},
		{
			Name:    "MultiLine/Bytes",
			Input:   kBytesValue,
			Factory: multiLineJSON,
			Expect:  ParseOneLine(`"YWJj"`),
		},
		{
			Name:    "MultiLine/Array",
			Input:   kArrayValue,
			Factory: multiLineJSON,
			Expect: ParseMultiLine(`
			|[
			|  "a",
			|  "b",
			|  "c"
			|]
			`),
		},
		{
			Name:    "MultiLine/Object",
			Input:   kObjectValue,
			Factory: multiLineJSON,
			Expect: ParseMultiLine(`
			|{
			|  "a": 1,
			|  "b": 2,
			|  "c": 3
			|}
			`),
		},
		{
			Name:    "MultiLine/Fancy",
			Input:   kFancyValue,
			Factory: multiLineJSON,
			Expect: ParseMultiLine(`
			|{
			|  "@type": "Foo",
			|  "emptyList": [],
			|  "emptyObject": {},
			|  "array": [
			|    "a",
			|    "b",
			|    "c"
			|  ],
			|  "object": {
			|    "a": 1,
			|    "b": 2,
			|    "c": 3
			|  }
			|}
			`),
		},
	}

	var buf bytes.Buffer
	var e emitter.Emitter
	for _, row := range testData {
		t.Run(row.Name, func(t *testing.T) {
			buf.Reset()
			e.Reset(&buf, row.Factory.NewGenerator())
			row.Input.EmitTo(&e)
			err := e.Close()

			if !reflect.DeepEqual(err, row.ExpectErr) {
				t.Errorf("wrong error:\n\texpect: %#v\n\tactual: %#v", row.ExpectErr, err)
				return
			}

			if result := buf.Bytes(); !bytes.Equal(result, row.Expect) {
				switch {
				case !IsSafe(result) || !IsSafe(row.Expect):
					expect := Hex(row.Expect)
					actual := Hex(result)
					t.Errorf("wrong result:\n\texpect: %#v\n\tactual: %#v", expect, actual)

				case IsMultiLine(result) || IsMultiLine(row.Expect):
					expect := FormatMultiLine(row.Expect, "\t|")
					actual := FormatMultiLine(result, "\t|")
					t.Errorf("wrong result:\n\texpect:\n%s\n\tactual:\n%s", expect, actual)

				default:
					expect := string(row.Expect)
					actual := string(result)
					t.Errorf("wrong result:\n\texpect: %q\n\tactual: %q", expect, actual)
				}
			}
		})
	}
}

func IsSafeByte(ch byte) bool {
	switch ch {
	case '\t':
		return true
	case '\n':
		return true
	case '\r':
		return true
	case 0x7f:
		return false
	default:
		return ch >= 0x20
	}
}

func IsSafe(buf []byte) bool {
	for _, ch := range buf {
		if !IsSafeByte(ch) {
			return false
		}
	}
	return true
}

func IsMultiLine(buf []byte) bool {
	n := uint(len(buf))
	if n <= 0 {
		return false
	}
	n--
	buf = buf[:n]
	for _, ch := range buf {
		if ch == '\n' {
			return true
		}
	}
	return false
}

func ParseOneLine(in string) []byte {
	out := make([]byte, 0, len(in)+1)
	out = append(out, in...)
	out = append(out, '\n')
	return out
}

func ParseMultiLine(in string) []byte {
	out := make([]byte, 0, len(in))
	active := false
	for _, ch := range in {
		switch {
		case active && ch == '\n':
			out = append(out, '\n')
			active = false
		case active:
			out = utf8.AppendRune(out, ch)
		case ch == '|':
			active = true
		}
	}
	return out
}

func FormatMultiLine(in []byte, prefix string) string {
	out := make([]byte, 0, len(in))
	needPrefix := true
	for _, ch := range in {
		switch {
		case ch == '\n':
			out = append(out, '\n')
			needPrefix = true
		case needPrefix:
			out = append(out, prefix...)
			needPrefix = false
			fallthrough
		default:
			out = append(out, ch)
		}
	}
	return string(out)
}

type Hex []byte

func (hex Hex) Append(out []byte) []byte {
	const DIGITS = "0123456789abcdef"
	out = append(out, '[')
	for index, ch := range hex {
		if index > 0 {
			out = append(out, ':')
		}
		hi := (ch >> 4) & 15
		lo := (ch & 15)
		out = append(out, DIGITS[hi], DIGITS[lo])
	}
	out = append(out, ']')
	return out
}

func (hex Hex) GoString() string {
	return string(hex.Append(nil))
}

func (hex Hex) String() string {
	return string(hex.Append(nil))
}
