package events

type GroupJoinRequestStarted struct {
	EventBase `json:"-"`
	Request   uint `json:"request"`
}

type GroupJoinRequestAccepted struct {
	EventBase `json:"-"`
	Request   uint `json:"request"`
}

type GroupJoinRequestDenied struct {
	EventBase `json:"-"`
	Request   uint `json:"request"`
}

type GroupJoinRequestWithdrawn struct {
	EventBase `json:"-"`
	Request   uint `json:"request"`
}
