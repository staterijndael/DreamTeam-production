package events

type OrgRenamed struct {
	EventBase `json:"-"`
	OldName   string `json:"oldName"`
	Org       uint   `json:"organization"`
	Director  uint   `json:"director"`
}

type OrgDeleted struct {
	EventBase    `json:"-"`
	Org          uint   `json:"organization"`
	GroupMembers []uint `json:"-"`
}

type UserAssociated struct {
	EventBase    `json:"-"`
	Org          uint `json:"organization"`
	Associated   uint `json:"associated"`
	AssociatedBy uint `json:"associatedBy"`
}

type UserDissociated struct {
	EventBase `json:"-"`
	Org       uint `json:"organization"`
	User      uint `json:"user"`
}

type UserDissociatedByDirector struct {
	EventBase `json:"-"`
	Org       uint `json:"organization"`
	User      uint `json:"user"`
	Director  uint `json:"director"`
}

type OrgNewDirector struct {
	EventBase   `json:"-"`
	Org         uint `json:"organization"`
	OldDirector uint `json:"oldDirector"`
	NewDirector uint `json:"newDirector"`
}
