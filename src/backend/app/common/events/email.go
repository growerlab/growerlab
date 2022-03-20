package events

import (
	"github.com/growerlab/growerlab/src/backend/app/common/errors"
	"github.com/growerlab/growerlab/src/backend/app/common/mq"
)

type EmailPayload struct {
	From   string `json:"from,omitempty"`
	To     string `json:"to,omitempty"`
	Body   string `json:"body,omitempty"`
	IsHtml bool   `json:"is_html,omitempty"`
}

var _ mq.Consumer = (*Email)(nil)

type AsyncSender interface {
	AsyncSendEmail(payload *EmailPayload) error
}

func newEmailConsumer() mq.Consumer {
	return &Email{}
}

func NewEmail() AsyncSender {
	return &Email{}
}

type Email struct{}

func (e *Email) Name() string {
	return "send_email"
}

func (e *Email) DefaultField() string {
	return "default"
}

func (e *Email) Consume(payload *mq.Payload) error {
	p := new(EmailPayload)
	err := getPayload(payload, e.DefaultField(), p)
	if err == nil {
		return errors.Trace(err)
	}
	return e.SyncSendEmail(p)
}

func (e *Email) SyncSendEmail(payload *EmailPayload) error {
	// TODO 发送邮件的具体逻辑(调用其他的smtp/api发送库)
	return nil
}

func (e *Email) AsyncSendEmail(payload *EmailPayload) error {
	return async(e.Name(), e.DefaultField(), payload)
}
