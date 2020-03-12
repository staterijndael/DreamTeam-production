package events

import "context"

type IEvent interface {
	GetContext() context.Context
}

type EventBase struct {
	Context context.Context `json:"-"`
}

func (e *EventBase) GetContext() context.Context {
	return e.Context
}
