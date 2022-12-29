package mq

import (
	"context"

	"github.com/growerlab/growerlab/src/common/jsonutils"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/growerlab/growerlab/src/common/db"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/common/logger"
	"github.com/minghsu0107/watermill-redistream/pkg/redis"
)

type MQ interface {
	RegisterHandler(uniqName, topic string, handlerFunc message.HandlerFunc) *message.Handler
	Publish(topic string, messages ...*message.Message) error
	DirectlyPublish(topic string, payload any) error
	Close() error
}

type RedisMQ struct {
	router *message.Router
	logger watermill.LoggerAdapter

	publisher  message.Publisher
	subscriber message.Subscriber
}

func (m *RedisMQ) DirectlyPublish(topic string, payload any) error {
	raw, err := jsonutils.EncodeObjectToBytes(payload)
	if err != nil {
		return errors.Trace(err)
	}
	msg := message.NewMessage(watermill.NewUUID(), raw)
	return m.Publish(topic, msg)
}

func NewRedisMQ() (MQ, error) {
	ctx := context.Background()
	lg := watermill.NewStdLoggerWithOut(logger.LogWriter, false, false)
	router, err := message.NewRouter(message.RouterConfig{}, lg)
	if err != nil {
		return nil, errors.Trace(err)
	}

	messagePublisher, err := redis.NewPublisher(ctx, redis.PublisherConfig{}, db.MemDB, redis.DefaultMarshaller{}, lg)
	if err != nil {
		return nil, errors.Trace(err)
	}

	messageSubscriber, err := redis.NewSubscriber(ctx, redis.SubscriberConfig{}, db.MemDB, redis.DefaultMarshaller{}, lg)
	if err != nil {
		return nil, errors.Trace(err)
	}

	instance := &RedisMQ{
		logger:     watermill.NewStdLoggerWithOut(logger.LogWriter, false, false),
		router:     router,
		publisher:  messagePublisher,
		subscriber: messageSubscriber,
	}

	go func() {
		err = router.Run(ctx)
		if err != nil {
			panic(err)
		}
	}()
	<-router.Running()
	return instance, nil
}

func (m *RedisMQ) RegisterHandler(uniqName, topic string, handlerFunc message.HandlerFunc) *message.Handler {
	return m.router.AddHandler(uniqName, topic, m.subscriber, "", m.publisher, handlerFunc)
}

func (m *RedisMQ) Publish(topic string, messages ...*message.Message) error {
	return m.publisher.Publish(topic, messages...)
}

func (m *RedisMQ) Close() error {
	return m.router.Close()
}
