package messaging

import (
	"context"
	"twitter-clone/internal/config"
	repositories "twitter-clone/internal/repositories/feed"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type NullMessageHandler struct {
}

func (n *NullMessageHandler) SetupMessageRouter(
	configuration config.Configuration,
	feedRepo repositories.FeedRepository,
	logger watermill.LoggerAdapter,
) (message.Publisher, message.Subscriber, error) {
	// Create no-op publisher and subscriber
	publisher := &NullPublisher{}
	subscriber := &NullSubscriber{}

	return publisher, subscriber, nil
}

// NullPublisher is a no-op publisher
type NullPublisher struct {
}

// Publish does nothing
func (p *NullPublisher) Publish(topic string, messages ...*message.Message) error {
	return nil
}

// Close does nothing
func (p *NullPublisher) Close() error {
	return nil
}

// NullSubscriber is a no-op subscriber
type NullSubscriber struct {
}

// Subscribe returns a closed channel
func (s *NullSubscriber) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	ch := make(chan *message.Message)
	close(ch)
	return ch, nil
}

// Close does nothing
func (s *NullSubscriber) Close() error {
	return nil
}
