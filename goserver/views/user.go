package views

import (
	"dt/models"
	"encoding/json"
)

type User struct {
	ID             uint          `json:"id"`
	Phone          *string       `json:"phone,omitempty"`
	Teel           uint64        `json:"teel"`
	Experience     uint64        `json:"experience"`
	Date           int64         `json:"date"`
	Score          Score         `json:"score"`
	Nickname       string        `json:"nickname,"`
	FirstName      *string       `json:"firstName,omitempty"`
	LastName       *string       `json:"lastName,omitempty"`
	AboutMe        *string       `json:"aboutMe,omitempty"`
	Email          *string       `json:"email,omitempty"`
	BDate          *int64        `json:"bdate,omitempty"`
	CreatedAt      int64         `json:"createdAt"`
	AvatarMetaInfo *FileMetaInfo `json:"avatarMetaInfo"`
}

type Score struct {
	Initiative int32  `json:"initiative"`
	Discipline int32  `json:"discipline"`
	Efficiency int32  `json:"efficiency"`
	Teamwork   int32  `json:"teamwork"`
	Loyalty    int32  `json:"loyalty"`
	Rang       string `json:"rang,omitempty"`
}

func ScoreViewFromModel(s *models.Score) Score {
	return Score{
		Initiative: s.Initiative,
		Discipline: s.Discipline,
		Efficiency: s.Efficiency,
		Teamwork:   s.Teamwork,
		Loyalty:    s.Loyalty,
		Rang:       s.Rang,
	}
}

func UserViewFromModel(u *models.User) *User {
	if u == nil {
		return nil
	}

	var score models.Score
	json.Unmarshal(u.Score.RawMessage, &score)
	var bDate *int64
	if u.BDate == nil {
		bDate = nil
	} else {
		_bDate := u.BDate.Unix()
		bDate = &_bDate
	}

	return &User{
		ID:             u.ID,
		Teel:           u.Teel,
		Experience:     u.Experience,
		Date:           u.CreatedAt.Unix(),
		Score:          ScoreViewFromModel(&score),
		Nickname:       u.Nickname.Value,
		FirstName:      stringOrNil(u.FirstName),
		LastName:       stringOrNil(u.LastName),
		AboutMe:        stringOrNil(u.AboutMe),
		Email:          stringOrNil(u.Email),
		BDate:          bDate,
		CreatedAt:      u.CreatedAt.Unix(),
		AvatarMetaInfo: FileMetaInfoViewFromModel(&u.Avatar),
	}
}

func UserSelfViewFromModel(u *models.User) *User {
	base := UserViewFromModel(u)
	base.Phone = &u.Phone
	return base
}
