package mq

import (
	"github.com/go-redis/redis/v7"
	"github.com/growerlab/growerlab/src/backend/app/common/errors"
)

var (
	ErrNoSuchKey = errors.New("ERR no such key")
)

type Stream struct {
	memDB redis.Cmdable
}

func NewStream(c redis.Cmdable) *Stream {
	return &Stream{memDB: c}
}

func (s *Stream) GroupExists(groupName string) bool {
	_, err := s.GroupInfo(groupName)
	if err != nil && err.Error() == ErrNoSuchKey.Error() {
		return false
	}
	return true
}

func (s *Stream) CreateGroup(groupName, streamKey string) error {
	if s.memDB.Exists(streamKey).Val() == 1 {
		return nil
	}
	err := s.memDB.XGroupCreateMkStream(streamKey, groupName, "0-0").Err()
	return errors.Trace(err)
}

func (s *Stream) GroupInfo(groupName string) ([]redis.XInfoGroups, error) {
	info, err := s.memDB.XInfoGroups(groupName).Result()
	return info, errors.Trace(err)
}

func (s *Stream) ReadGroupMessages(groupName, consumer, streamKey string, count int64) ([]redis.XMessage, error) {
	msgs, err := s.readGroupMessages(groupName, consumer, []string{streamKey, ">"}, count)
	if err != nil {
		return nil, err
	}
	if len(msgs) > 0 {
		return msgs, nil
	}

	// 读取历史数据
	msgs, err = s.readGroupMessages(groupName, consumer, []string{streamKey, "0-0"}, count)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func (s *Stream) readGroupMessages(groupName, consumer string, streams []string, count int64) ([]redis.XMessage, error) {
	xstreams, err := s.memDB.XReadGroup(&redis.XReadGroupArgs{
		Group:    groupName,
		Consumer: consumer,
		Streams:  streams,
		Count:    count,
		Block:    0,
		NoAck:    true,
	}).Result()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if len(xstreams) == 0 {
		return nil, nil
	}

	return xstreams[0].Messages, nil
}

func (s *Stream) AddMessage(streamKey, field, value string) (id string, err error) {
	id, err = s.memDB.XAdd(&redis.XAddArgs{
		Stream: streamKey,
		ID:     "*",
		Values: map[string]interface{}{field: value},
	}).Result()

	return id, errors.Trace(err)
}

func (s *Stream) Ack(streamKey, groupName, id string) (n int64, err error) {
	n, err = s.memDB.XAck(streamKey, groupName, id).Result()
	return n, errors.Trace(err)
}
