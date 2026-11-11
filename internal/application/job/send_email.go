package job

import (
	"maps"
	"time"
)

var _ Job = (*SendEmail)(nil)

type SendEmail struct {
	BaseJob

	to         string
	templateID string
	variables  map[string]string
}

func NewSendEmail(
	to string,
	templateID string,
	variables map[string]string,
	enqueueTime time.Time,
) SendEmail {
	return SendEmail{
		BaseJob: NewBaseJob(enqueueTime),

		to:         to,
		templateID: templateID,
		variables:  maps.Clone(variables),
	}
}

func (e SendEmail) Type() string {
	return "email.send-email"
}

func (e SendEmail) To() string {
	return e.to
}

func (e SendEmail) TemplateID() string {
	return e.templateID
}

func (e SendEmail) Variables() map[string]string {
	return maps.Clone(e.variables)
}
