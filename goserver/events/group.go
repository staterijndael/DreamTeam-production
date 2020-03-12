package events

type GroupCreated struct {
	EventBase `json:"-"`
	Group     uint `json:"group"`
	Creator   uint `json:"creator"`
}

type GroupsDeleted struct {
	EventBase `json:"-"`
	Groups    []uint `json:"groups"`
	DeletedBy uint   `json:"deletedBy"`
}

type GroupRenamed struct {
	EventBase `json:"-"`
	OldName   string `json:"oldName"`
	Group     uint   `json:"group"`
	RenamedBy uint   `json:"renamedBy"`
}

type GroupAdminChanged struct {
	EventBase `json:"-"`
	Group     uint `json:"group"`
	NewAdmin  uint `json:"newAdmin"`
	OldAdmin  uint `json:"oldAdmin"`
	ChangedBy uint `json:"changedBy"`
}

type UserAddedToGroup struct {
	EventBase `json:"-"`
	Group     uint `json:"group"`
	Added     uint `json:"added"`
	AddedBy   uint `json:"addedBy"`
}

type UserLeftGroup struct {
	EventBase `json:"-"`
	Group     uint `json:"group"`
	User      uint `json:"user"`
}

type UserRemovedFromGroup struct {
	EventBase `json:"-"`
	Group     uint `json:"group"`
	Removed   uint `json:"removed"`
	RemovedBy uint `json:"removedBy"`
}

type NewGroupAdmin struct {
	EventBase  `json:"-"`
	Group      uint `json:"group"`
	NewAdmin   uint `json:"newAdmin"`
	OldAdmin   uint `json:"oldAdmin"`
	AssignedBy uint `json:"assignedBy"`
}

type GroupDeleted struct {
	EventBase `json:"-"`
	Group     uint `json:"group"`
	DeletedBy uint `json:"deletedBy"`
}
