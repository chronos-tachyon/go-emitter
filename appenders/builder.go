package appenders

type Builder struct {
	list []Appender
	out  []byte
}

func (b *Builder) flush() {
	if len(b.out) <= 0 {
		return
	}
	str := string(b.out)
	b.out = nil
	b.list = append(b.list, LiteralString(str))
}

func (b *Builder) Add(a Appender) {
	switch x := a.(type) {
	case nil:
		// pass
	case LiteralBytes:
		b.AddBytes([]byte(x))
	case LiteralString:
		b.AddString(string(x))
	default:
		b.flush()
		b.list = append(b.list, a)
	}
}

func (b *Builder) Indent(tab bool, size uint, count uint) {
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

func (b *Builder) AddBytes(buf []byte) {
	b.out = append(b.out, buf...)
}

func (b *Builder) AddString(str string) {
	b.out = append(b.out, str...)
}

func (b *Builder) AddByte(ch byte) {
	b.out = append(b.out, ch)
}

func (b *Builder) AddRune(ch rune) {
	str := string(ch)
	b.out = append(b.out, str...)
}

func (b *Builder) Build() []Appender {
	b.flush()
	list := b.list
	b.list = nil
	return list
}
