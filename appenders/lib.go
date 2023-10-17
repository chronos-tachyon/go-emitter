package appenders

type Appender interface {
	Append(out []byte) []byte
}
