package events

import (
	"dt/models"
	"encoding/json"
)

type UserRenamed struct {
	EventBase `json:"-"`
	OldUser   models.User
	User      models.User
}

type UserAccountDeleted struct {
	EventBase `json:"-"`
	User      uint   `json:"user"`
	Groups    []uint `json:"-"`
}

func (event *UserRenamed) MarshallJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m["oldName"] = event.OldUser.FirstName.String + " " + event.OldUser.LastName.String
	m["newName"] = event.User.FirstName.String + " " + event.User.LastName.String
	return json.Marshal(m)
}
