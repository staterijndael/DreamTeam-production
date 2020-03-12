package models

import (
	"github.com/jinzhu/gorm"
)

type RequestStatus string

const (
	Denied    RequestStatus = "denied"
	Confirmed RequestStatus = "confirmed"
	Pending   RequestStatus = "pending"
	Canceled  RequestStatus = "canceled"
)

type RequestBase struct {
	Status      RequestStatus `gorm:"column:status"`
	InitiatorID uint          `gorm:"column:initiator"`
	AcceptorID  *uint         `gorm:"column:acceptor"`
	Initiator   User          `gorm:"foreignkey:initiator"`
	Acceptor    *User         `gorm:"foreignkey:acceptor"`
}

//TODO rename parent column name to parent_id and test it
type GroupCreationRequest struct {
	gorm.Model
	RequestBase
	OrganizationID uint         `gorm:"column:organization"`
	ParentID       *uint        `gorm:"column:hm"`
	Title          string       `gorm:"column:title"`
	Description    string       `gorm:"column:description"`
	NicknameID     uint         `gorm:"column:nickname"`
	Nickname       Nickname     `gorm:"foreignkey:nickname"`
	Parent         *Group       `gorm:"foreignkey:hm"`
	Organization   Organization `gorm:"foreignkey:organization"`
}

func (*GroupCreationRequest) TableName() string {
	return "group_creation_requests"
}

type GroupJoinRequest struct {
	gorm.Model
	RequestBase
	GroupID uint  `gorm:"column:group"`
	Group   Group `gorm:"foreignkey:group"`
}

type OrgJoinRequest struct {
	gorm.Model
	RequestBase
	OrganizationID uint         `gorm:"column:organization"`
	GroupID        *uint        `gorm:"column:group"`
	Organization   Organization `gorm:"foreignkey:organization"`
	Group          *Group       `gorm:"foreignkey:group"`
}

func (*OrgJoinRequest) TableName() string {
	return "org_join_requests"
}

func (*GroupJoinRequest) TableName() string {
	return "group_join_requests"
}
