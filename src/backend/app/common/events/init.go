package events

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/growerlab/growerlab/src/backend/app/common/mq"
	"github.com/growerlab/growerlab/src/backend/app/common/notify"
	"github.com/growerlab/growerlab/src/common/logger"
)

var eventMQ mq.MQ

type EventProcessor interface {
	Topic() string
	Handler(msg *message.Message) ([]*message.Message, error)
}

type EventPublisher interface {
	Publish(topic string, messages ...*message.Message) error
}

func InitEvents() (err error) {
	eventMQ, err = mq.NewRedisMQ()

	notify.Subscribe(func() {
		err = eventMQ.Close()
		if err != nil {
			logger.ErrorTrace(err)
		}
	})

	events := []EventProcessor{
		NewEmailProcessor(),
		NewGitEventProcessor(),
	}

	for _, event := range events {
		uniqName := "mq:" + event.Topic()
		_ = eventMQ.RegisterHandler(uniqName, event.Topic(), event.Handler)
		if err != nil {
			return err
		}
	}

	return
}
