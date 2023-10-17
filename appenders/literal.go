package appenders

type LiteralBytes []byte

func (a LiteralBytes) Append(out []byte) []byte {
	slice := []byte(a)
	return append(out, slice...)
}

var _ Appender = LiteralBytes(nil)

type LiteralString string

func (a LiteralString) Append(out []byte) []byte {
	str := string(a)
	return append(out, str...)
}

var _ Appender = LiteralString("")
