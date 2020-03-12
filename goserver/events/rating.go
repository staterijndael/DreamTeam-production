package events

import (
	"encoding/json"
	"errors"
	"time"
)

type UserEstimated struct {
	User uint `json:"user"`
}

type RatingStartedEvent struct {
	Start          time.Time
	End            time.Time
	OrganizationID uint
}

type RatingEndedEvent struct {
	End            time.Time
	OrganizationId uint
}

type RatingEnabled struct {
	EventBase    `json:"-"`
	Organization uint `json:"organization"`
	Config       uint `json:"config"`
}

type RatingDisabled struct {
	EventBase    `json:"-"`
	Organization uint `json:"organization"`
}

func (self *RatingStartedEvent) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})

	m["start"] = self.Start.Unix()
	m["end"] = self.End.Unix()
	m["organization"] = self.OrganizationID

	return json.Marshal(m)
}

func (self *RatingStartedEvent) UnmarshalJSON(data []byte) error {
	m := make(map[string]interface{})

	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	start, ok1 := m["start"]
	end, ok2 := m["end"]
	org, ok3 := m["organization"]

	if !ok1 || !ok2 || !ok3 {
		return errors.New("invalid JSON for RatingStartedEvent")
	}

	startT, ok1 := start.(int64)
	endT, ok2 := end.(int64)
	orgId, ok3 := org.(uint)

	if !ok1 || !ok2 || !ok3 {
		return errors.New("invalid JSON for RatingStartedEvent")
	}

	self.Start = time.Unix(startT, 0)
	self.End = time.Unix(endT, 0)
	self.OrganizationID = orgId

	return nil
}
