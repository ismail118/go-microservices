package event

type Emitter interface {
	SetupEvenEmitter() error
	Push(event string, severity string) error
}
