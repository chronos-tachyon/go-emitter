package emitter

import (
	"encoding"
	"fmt"
)

type JSONFormat byte

const (
	CompactJSON JSONFormat = iota
	OneLineJSON
	MultiLineJSON
)

const jsonFormatSize = 3

var jsonFormatGoNames = [jsonFormatSize]string{
	"emitter.CompactJSON",
	"emitter.OneLineJSON",
	"emitter.MultiLineJSON",
}

var jsonFormatNames = [jsonFormatSize]string{
	"compact",
	"oneLine",
	"multiLine",
}

func (f JSONFormat) IsValid() bool {
	return f < jsonFormatSize
}

func (f JSONFormat) GoString() string {
	if f.IsValid() {
		return jsonFormatGoNames[f]
	}
	return fmt.Sprintf("emitter.JSONFormat(%d)", uint(f))
}

func (f JSONFormat) String() string {
	if f.IsValid() {
		return jsonFormatNames[f]
	}
	return fmt.Sprintf("%%!ERR[invalid JSONFormat %d]", uint(f))
}

func (f JSONFormat) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

func (f *JSONFormat) Parse(input string) error {
	*f = ^JSONFormat(0)
	return fmt.Errorf("failed to parse %q as emitter.JSONFormat", input)
}

func (f *JSONFormat) UnmarshalText(input []byte) error {
	return f.Parse(string(input))
}

func (f JSONFormat) indent(b *AppenderBuilder, tabs bool, size uint, count uint) {
	switch f {
	case MultiLineJSON:
		b.AddByte('\n')
		b.Indent(tabs, size, count)
	}
}

func (f JSONFormat) indentOrSpace(b *AppenderBuilder, tabs bool, size uint, count uint) {
	switch f {
	case MultiLineJSON:
		b.AddByte('\n')
		b.Indent(tabs, size, count)
	case OneLineJSON:
		b.AddByte(' ')
	}
}

func (f JSONFormat) space(b *AppenderBuilder) {
	switch f {
	case MultiLineJSON:
		fallthrough
	case OneLineJSON:
		b.AddByte(' ')
	}
}

func (f JSONFormat) lineFeed(b *AppenderBuilder) {
	switch f {
	case MultiLineJSON:
		fallthrough
	case OneLineJSON:
		b.AddByte('\n')
	}
}

var (
	_ fmt.GoStringer           = JSONFormat(0)
	_ fmt.Stringer             = JSONFormat(0)
	_ encoding.TextMarshaler   = JSONFormat(0)
	_ encoding.TextUnmarshaler = (*JSONFormat)(nil)
)
