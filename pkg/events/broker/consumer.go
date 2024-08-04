package broker

type Consumer interface {
	StartConsumer(processMessage func(message []byte) error)
}
