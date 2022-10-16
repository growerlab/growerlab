package mq

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/growerlab/growerlab/src/backend/app/utils/logger"
	"github.com/growerlab/growerlab/src/common/db"
	"github.com/ivpusic/grpool"
)

const (
	DefaultField     = "default"
	DefaultValue     = "default"
	DefaultGroup     = "defaultGroup"
	DefaultConsumer  = "GrowerLab"
	DefaultReadCount = 10
)

const (
	JobWorkers = 32
	JobQueue   = 32
)

type Payload struct {
	Consumer string // 所属消息ID
	ID       string
	Values   map[string]interface{}
}

func (p *Payload) Get(fd string) interface{} {
	if v, ok := p.Values[fd]; ok {
		return v
	}
	return nil
}

type Consumer interface {
	Name() string           // consumer name
	Consume(*Payload) error // 进行消费
	// 下面的功能先注释，未来再加
	// Number() int            // 消费者人数（决定提供多少worker调用Consume()
	// RetryCount() int        // 重试次数
	// RetryInterval() int     // 重试间隔，单位s
}

type MessageQueue struct {
	memDB  *db.MemDBClient
	stream *Stream

	consumers       map[string]Consumer // 消费者
	waitingMessages chan *Payload       // 等待处理的消息

	pool        *grpool.Pool
	done        chan struct{}
	releaseOnce sync.Once
}

func NewMessageQueue(c *db.MemDBClient) *MessageQueue {
	return &MessageQueue{
		memDB:           c,
		stream:          NewStream(c),
		waitingMessages: make(chan *Payload, 1024),
		pool:            grpool.NewPool(JobWorkers, JobQueue),
		done:            make(chan struct{}),
	}
}

func (m *MessageQueue) Register(consumers ...Consumer) error {
	if len(consumers) == 0 {
		return nil
	}
	if m.consumers == nil {
		m.consumers = map[string]Consumer{}
	}

	for _, c := range consumers {
		if _, exists := m.consumers[c.Name()]; exists {
			return fmt.Errorf("consumer exists: %s", c.Name())
		}
		m.consumers[c.Name()] = c
		err := m.createStream(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MessageQueue) createStream(c Consumer) error {
	if exists := m.stream.GroupExists(DefaultGroup); exists {
		return nil
	}

	streamKey := m.streamKey(c.Name())

	return m.stream.CreateGroup(DefaultGroup, streamKey)
}

func (m *MessageQueue) streamKey(name string) string {
	return m.memDB.KeyBuilder.KeyMaker().Append(name).String()
}

func (m *MessageQueue) buildPayload(belongID string, msg *redis.XMessage) *Payload {
	payload := &Payload{
		Consumer: belongID,
		ID:       msg.ID,
		Values:   msg.Values,
	}
	return payload
}

func (m *MessageQueue) Run() error {
	go func() {
		for {
			select {
			case <-m.done:
				return
			default:
				count := m.take()

				// 延迟的目的是降低内存数据库的压力
				if count == 0 {
					time.Sleep(1 * time.Second)
				} else {
					time.Sleep(50 * time.Millisecond)
				}
			}
		}
	}()
	go func() {
		for {
			select {
			case <-m.done:
				return
			case msg := <-m.waitingMessages:
				m.delivery(msg)
			}
		}
	}()
	return nil
}

func (m *MessageQueue) delivery(msg *Payload) {
	defer func() {
		if e := recover(); e != nil {
			logger.Error("[message queue] delivery err: %v", e)
		}
	}()
	consumer, ok := m.consumers[msg.Consumer]
	if !ok {
		logger.Error("[message queue] not found consumer: %s", msg.Consumer)
		return
	}
	m.pool.JobQueue <- func() {
		err := consumer.Consume(msg)
		if err != nil {
			logger.Error("[message queue] consumer consume was err: %v", err)
			return
		} else { // err == nil
			streamKey := m.streamKey(msg.Consumer)
			_, err = m.stream.Ack(streamKey, DefaultGroup, msg.ID)
			if err != nil {
				logger.Error("[message queue] ack was err: %s - %s - %v", streamKey, msg.ID, err)
				return
			}
		}
	}
}

func (m *MessageQueue) take() (count int) {
	defer func() {
		if e := recover(); e != nil {
			logger.Error("[message queue] running err: %v", e)
		}
	}()
	for _, c := range m.consumers {
		streamKey := m.streamKey(c.Name())
		messages, err := m.stream.ReadGroupMessages(DefaultGroup, DefaultConsumer, streamKey, DefaultReadCount)
		if err != nil {
			panic(err)
		}
		count += len(messages)

		for _, msg := range messages {
			m.waitingMessages <- m.buildPayload(c.Name(), &msg)
		}
	}
	return
}

func (m *MessageQueue) Add(consumerName, msgField, msgBody string) (id string, err error) {
	_, ok := m.consumers[consumerName]
	if !ok {
		return "", fmt.Errorf("not found consumer: %s", consumerName)
	}

	streamKey := m.streamKey(consumerName)
	return m.stream.AddMessage(streamKey, msgField, msgBody)
}

func (m *MessageQueue) Release() {
	m.releaseOnce.Do(func() {
		close(m.done)
		m.pool.Release()
	})
}
