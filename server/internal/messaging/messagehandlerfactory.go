package messaging

import (
	"errors"
	"twitter-clone/internal/config"
)

func CreateMessageHandler(configuration config.Configuration) (MessageHandler, error) {
	switch configuration.Mode {
	case config.InMemory:
		return nil, errors.ErrUnsupported
	case config.Persistent:
		return &NATSMessageHandler{}, nil
	case config.Cloud:
		return &PubSubMessageHandler{}, nil
	default:
		return nil, errors.New("unknown mode")
	}
}
