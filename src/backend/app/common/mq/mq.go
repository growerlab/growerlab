package mq

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/growerlab/growerlab/src/common/db"
	"github.com/minghsu0107/watermill-redistream/pkg/redis"
)

var logger = watermill.NewStdLogger(false, false)
var router *message.Router

var (
	messagePublisher  message.Publisher
	messageSubscriber message.Subscriber
)

func init() {
	var err error
	router, err = message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	messagePublisher, err = redis.NewPublisher(ctx, redis.PublisherConfig{}, db.MemDB, redis.DefaultMarshaller{}, logger)
	if err != nil {
		panic(err)
	}

	messageSubscriber, err = redis.NewSubscriber(ctx, redis.SubscriberConfig{}, db.MemDB, redis.DefaultMarshaller{}, logger)
	if err != nil {
		panic(err)
	}
}

func RegisterEvents(uniqName, topic string, handlerFunc message.HandlerFunc) *message.Handler {
	return router.AddHandler(uniqName, topic, messageSubscriber, "", messagePublisher, handlerFunc)
}
