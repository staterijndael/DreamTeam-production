package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"math"
	"strings"
	"time"
)

type User struct {
	gorm.Model
	AvatarID   uint           `gorm:"column:avatar"`
	Teel       uint64         `gorm:"column:teel"`
	Phone      string         `gorm:"column:phone; unique"`
	Password   string         `gorm:"column:password"`
	Experience uint64         `gorm:"column:experience"`
	BDate      *time.Time     `gorm:"column:bdate"`
	Score      postgres.Jsonb `gorm:"column:score"`
	NicknameID *uint          `gorm:"column:nickname"`
	Nickname   *Nickname      `gorm:"foreignkey:nickname"`
	FirstName  sql.NullString `gorm:"column:first_name"`
	LastName   sql.NullString `gorm:"column:last_name"`
	AboutMe    sql.NullString `gorm:"column:about_me"`
	Email      sql.NullString `gorm:"column:email"`
	Avatar     File           `gorm:"foreignkey:avatar"`
}

func (*User) FuzzyQuery(queries []string) string {
	var b strings.Builder
	var fuzziness uint32
	fmt.Fprintf(&b, "%s ==> dsl.or(", new(User).TableName())
	for i := range queries {
		if len(queries[i]) > 3 {
			fuzziness = 2
		} else {
			fuzziness = 1
		}

		fmt.Fprintf(
			&b,
			"dsl.fuzzy('zdb_all', '%s', fuzziness => %d), dsl.wildcard('zdb_all', '*%s*'), ",
			queries[i],
			fuzziness,
			queries[i],
		)
	}

	return b.String()[:b.Len()-2] + ")"
}

var EmptyNicknameModelErr = errors.New("no nickname")

func (u *User) TableName() string {
	return "users"
}

func (u *User) ZDBIdxTypeDefinition() string {
	return "id " + u.ZDBIdxIDType() + ", first_name text, last_name text"
}

func (u *User) ZDBRowBuilder() string {
	return "(id)::" + u.ZDBIdxIDType() + ", first_name, last_name"
}

func (u *User) ZDBIdxIDType() string {
	return "bigint"
}

func (u *User) GetScore() *Score {
	var score Score
	json.Unmarshal(u.Score.RawMessage, &score)
	return &score
}

func (u *User) SetScore(s *Score) {
	bytes, _ := json.Marshal(s)
	u.Score = postgres.Jsonb{RawMessage: bytes}
}

var StartedScore json.RawMessage

func init() {
	b, _ := json.Marshal(GetStartedScore())
	StartedScore = b
}

func GetStartedScore() Score {
	s := Score{
		Initiative:    3,
		Discipline:    3,
		Efficiency:    3,
		Teamwork:      3,
		Loyalty:       3,
		AmountOfVotes: 1,
	}

	s.SetUpRang()
	return s
}

type Score struct {
	Initiative    int32  `json:"initiative"`
	Discipline    int32  `json:"discipline"`
	Efficiency    int32  `json:"efficiency"`
	Teamwork      int32  `json:"teamwork"`
	Loyalty       int32  `json:"loyalty"`
	Rang          string `json:"rang,omitempty"`
	AmountOfVotes uint64 `json:"amountOfVotes"`
}

func (s *Score) SetUpRang() {
	s.Rang = rangMap[int32(math.Round(float64(
		s.Initiative+s.Discipline+s.Efficiency+s.Teamwork+s.Loyalty)/5))]
}

var (
	rangMap = map[int32]string{
		0: "F",
		1: "E",
		2: "D",
		3: "C",
		4: "B",
		5: "A",
	}
)

func (s *Score) Adjust(score *Score) {

	amountF := float64(s.AmountOfVotes)

	curInitF := float64(s.Initiative)
	addInitF := float64(score.Initiative)
	s.Initiative = int32(math.Round((curInitF + addInitF/amountF) * (amountF / (amountF + 1))))

	curDiscF := float64(s.Discipline)
	addDiscF := float64(score.Discipline)
	s.Discipline = int32(math.Round((curDiscF + addDiscF/amountF) * (amountF / (amountF + 1))))

	curEffF := float64(s.Efficiency)
	addEffF := float64(score.Efficiency)
	s.Efficiency = int32(math.Round((curEffF + addEffF/amountF) * (amountF / (amountF + 1))))

	curTeamworkF := float64(s.Teamwork)
	addTeamworkF := float64(score.Teamwork)
	s.Teamwork = int32(math.Round((curTeamworkF + addTeamworkF/amountF) * (amountF / (amountF + 1))))

	curLoyF := float64(s.Loyalty)
	addLoyF := float64(score.Loyalty)
	s.Loyalty = int32(math.Round((curLoyF + addLoyF/amountF) * (amountF / (amountF + 1))))

	s.AmountOfVotes++
	s.SetUpRang()
}

func (s *Score) IsValid() bool {
	return s.Initiative >= 0 && s.Initiative <= 5 &&
		s.Discipline >= 0 && s.Discipline <= 5 &&
		s.Efficiency >= 0 && s.Efficiency <= 5 &&
		s.Teamwork >= 0 && s.Teamwork <= 5 &&
		s.Loyalty >= 0 && s.Loyalty <= 5
}
