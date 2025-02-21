package nats

import (
	"time"

	"github.com/nats-io/nats.go"
)

func Connect(url string) (*nats.Conn, error) {
	opts := []nats.Option{
		nats.Name("agent"),
		nats.ReconnectWait(5 * time.Second),
		nats.MaxReconnects(-1),
	}

	nc, err := nats.Connect(url, opts...)
	if err != nil {
		return nil, err
	}

	return nc, nil
}
