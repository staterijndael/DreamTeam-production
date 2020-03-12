package events

type GroupCreationRequestStarted struct {
	EventBase `json:"-"`
	Request   uint `json:"request"`
}

type GroupCreationRequestAccepted struct {
	EventBase `json:"-"`
	Request   uint  `json:"request"`
	Group     *uint `json:"group"`
}

type GroupCreationRequestDenied struct {
	EventBase `json:"-"`
	Request   uint `json:"request"`
}

type GroupCreationRequestWithdrawn struct {
	EventBase `json:"-"`
	Request   uint `json:"request"`
}
