package emitter

type Value interface {
	EmitTo(*Emitter)
}
