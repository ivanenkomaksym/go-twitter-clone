package messaging

import (
	"twitter-clone/internal/config"
	repositories "twitter-clone/internal/repositories/feed"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type MessageHandler interface {
	SetupMessageRouter(
		configuration config.Configuration,
		feedRepo repositories.FeedRepository,
		logger watermill.LoggerAdapter,
	) (message.Publisher, message.Subscriber, error)
}
