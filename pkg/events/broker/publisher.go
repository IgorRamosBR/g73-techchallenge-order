package broker

import "context"

type Publisher interface {
	Publish(ctx context.Context, destination string, message []byte) error
	Close() error
}
