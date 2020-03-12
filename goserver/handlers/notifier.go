package handlers

import (
	"dt/events"
	"dt/logwrap"
	"dt/notification"
	"dt/requestContext"
	"dt/views"
	"github.com/golang-collections/collections/set"
	"github.com/gorilla/websocket"
	"reflect"
)

func notifier(event interface{}) {
	handleError := func(er error) {
		value := reflect.ValueOf(event)
		if reflect.Ptr == value.Kind() {
			value = value.Elem()
		}

		logwrap.Error(
			"handler. event: %s{%v}; error: %s", value.Type().Name(), event, er.Error(),
		)
	}

	n, err := notification.GetNotification(event)
	if err != nil {
		handleError(err)
		return
	}

	if err = n.CreateByEvent(db, event); err != nil {
		handleError(err)
		return
	}

	var parentConnection *websocket.Conn = nil
	if converted, ok := event.(events.IEvent); ok {
		parentConnection, ok = requestContext.WebsocketFromContext(converted.GetContext())
		if !ok {
			parentConnection = nil
		}
	}

	msg := &views.JSONRPCNotification{
		Method: "new",
		Params: n.ContainerizedView(),
	}

	plebeianSet := receiversToSet(n.Receivers())
	for _, err := range sendToAllPlebeianMembers(receiversFromSet(plebeianSet), msg, parentConnection) {
		handleError(err)
	}

	dashSet := receiversToSet(n.DashReceivers()).Union(plebeianSet)
	for _, err := range sendToAllOrgAdminMembers(receiversFromSet(dashSet), msg, parentConnection) {
		handleError(err)
	}

	return
}

func receiversToSet(r []uint) *set.Set {
	s := set.New()
	for _, id := range r {
		s.Insert(id)
	}

	return s
}

func receiversFromSet(s *set.Set) []uint {
	a := make([]uint, 0, s.Len())
	s.Do(func(el interface{}) {
		a = append(a, el.(uint))
	})

	return a
}

func sendToAllPlebeianMembers(members []uint, msg interface{}, parent *websocket.Conn) []error {
	return sendToAllMembers(members, msg, parent, cm)
}

func sendToAllOrgAdminMembers(members []uint, msg interface{}, parent *websocket.Conn) []error {
	return sendToAllMembers(members, msg, parent, dcm)
}
