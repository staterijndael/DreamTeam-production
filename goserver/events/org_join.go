package events

type OrgJoinRequestStarted struct {
	EventBase `json:"-"`
	Request   uint `json:"request"`
}

type OrgJoinRequestAccepted struct {
	EventBase `json:"-"`
	Request   uint `json:"request"`
}

type OrgJoinRequestDenied struct {
	EventBase `json:"-"`
	Request   uint `json:"request"`
}

type OrgJoinRequestWithdrawn struct {
	EventBase `json:"-"`
	Request   uint `json:"request"`
}
