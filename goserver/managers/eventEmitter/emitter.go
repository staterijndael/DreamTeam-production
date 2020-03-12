package eventEmitter

import (
	"dt/config"
	"dt/handlers"
)

type HandlersStore []handlers.Handler

type EventEmitter struct {
	handlers HandlersStore
	events   chan interface{}
}

func NewEventEmitter(conf *config.Config, handlers []handlers.Handler) *EventEmitter {
	ee := &EventEmitter{
		handlers: handlers,
		events:   make(chan interface{}, conf.EventEmitterChannelBufferSize),
	}

	go ee.eventLoop()
	return ee
}

func (ee *EventEmitter) Listen(h handlers.Handler) {
	ee.handlers = append(ee.handlers, h)
}

func (ee *EventEmitter) Emit(event interface{}) {
	ee.events <- event
}

func (ee *EventEmitter) eventLoop() {
	for {
		eventWrapper := <-ee.events
		for _, h := range ee.handlers {
			h(eventWrapper)
		}
	}
}
