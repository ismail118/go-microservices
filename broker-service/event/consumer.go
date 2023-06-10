package event

type Consumer interface {
	Listen(topics []string) error
}
