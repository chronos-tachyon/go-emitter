package json

import (
	"github.com/chronos-tachyon/go-emitter"
)

type JSON struct {
	Format         Format
	IndentSize     uint
	IndentWithTabs bool
	EscapeHTML     bool
	TraceEnabled   bool
}

func (json JSON) NewGenerator() emitter.Generator {
	g := &Generator{json: json}
	g.Reset()
	return g
}

var _ emitter.GeneratorFactory = JSON{}
