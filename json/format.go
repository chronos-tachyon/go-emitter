package json

import (
	"encoding"
	"fmt"

	"github.com/chronos-tachyon/go-emitter/appenders"
)

type Format byte

const (
	Compact Format = iota
	OneLine
	MultiLine
)

const formatSize = 3

var formatGoNames = [formatSize]string{
	"json.Compact",
	"json.OneLine",
	"json.MultiLine",
}

var formatNames = [formatSize]string{
	"compact",
	"oneLine",
	"multiLine",
}

func (f Format) IsValid() bool {
	return f < formatSize
}

func (f Format) GoString() string {
	if f.IsValid() {
		return formatGoNames[f]
	}
	return fmt.Sprintf("json.Format(%d)", uint(f))
}

func (f Format) String() string {
	if f.IsValid() {
		return formatNames[f]
	}
	return fmt.Sprintf("%%!ERR[invalid json.Format %d]", uint(f))
}

func (f Format) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

func (f *Format) Parse(input string) error {
	*f = ^Format(0)
	return fmt.Errorf("failed to parse %q as json.Format", input)
}

func (f *Format) UnmarshalText(input []byte) error {
	return f.Parse(string(input))
}

func (f Format) indent(b *appenders.Builder, tabs bool, size uint, count uint) {
	switch f {
	case MultiLine:
		b.AddByte('\n')
		b.Indent(tabs, size, count)
	}
}

func (f Format) indentOrSpace(b *appenders.Builder, tabs bool, size uint, count uint) {
	switch f {
	case MultiLine:
		b.AddByte('\n')
		b.Indent(tabs, size, count)
	case OneLine:
		b.AddByte(' ')
	}
}

func (f Format) space(b *appenders.Builder) {
	switch f {
	case MultiLine:
		fallthrough
	case OneLine:
		b.AddByte(' ')
	}
}

func (f Format) lineFeed(b *appenders.Builder) {
	switch f {
	case MultiLine:
		fallthrough
	case OneLine:
		b.AddByte('\n')
	}
}

var (
	_ fmt.GoStringer           = Format(0)
	_ fmt.Stringer             = Format(0)
	_ encoding.TextMarshaler   = Format(0)
	_ encoding.TextUnmarshaler = (*Format)(nil)
)
