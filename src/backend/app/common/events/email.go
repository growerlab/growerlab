package events

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/growerlab/growerlab/src/common"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/common/jsonutils"
)

type EmailPayload struct {
	From   string `json:"from" validate:"required,email"`
	To     string `json:"to" validate:"required,email"`
	Body   string `json:"body" validate:"required,min=1"`
	IsHtml bool   `json:"is_html" validate:"required"`
}

var _ EventProcessor = (*Email)(nil)
var _ sender = (*Email)(nil)

type sender interface {
	PublishEmail(payload *EmailPayload) error
}

func NewEmailProcessor() EventProcessor {
	return &Email{}
}

func NewEmailSender() sender {
	return &Email{}
}

type Email struct{}

func (e *Email) PublishEmail(payload *EmailPayload) error {
	// 推送到 topic
	if err := common.Validator(payload); err != nil {
		return errors.Trace(err)
	}

	err := eventMQ.DirectlyPublish(e.Topic(), payload)
	return errors.Trace(err)
}

func (e *Email) Topic() string {
	return "send_email"
}

func (e *Email) Handler(msg *message.Message) ([]*message.Message, error) {
	p := new(EmailPayload)
	err := jsonutils.DecodeBytesToObject(msg.Payload, p)
	if err == nil {
		return nil, errors.Trace(err)
	}
	return nil, e.syncSendEmail(p)
}

func (e *Email) syncSendEmail(payload *EmailPayload) error {
	// TODO 发送邮件的具体逻辑(调用其他的smtp/api发送库)
	return nil
}
