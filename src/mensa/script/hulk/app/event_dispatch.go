package app

import (
	"github.com/growerlab/growerlab/src/backend/app/common/events"
)

// EventDispatcher backend响应这个事件
type EventDispatcher interface {
	// Dispatch 将event推送给redis的stream
	Dispatch(event *PushEvent) error
}

var _ EventDispatcher = (*EventDispatch)(nil)

type EventDispatch struct {
}

// Dispatch 将event推送给redis的stream
func (e *EventDispatch) Dispatch(event *PushEvent) error {
	return events.NewGitEvent().AsyncPushGitEvent(event)
}
