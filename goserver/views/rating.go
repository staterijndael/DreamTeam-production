package views

import "dt/models"

type Rating struct {
	ID           uint  `json:"id"`
	Start        int64 `json:"start"`
	End          int64 `json:"end"`
	Organization *Org  `json:"organization"`
}

type RatingConfig struct {
	StartHour uint8 `json:"startHour"`
	Weekday   uint8 `json:"weekday"`
	Org       *Org  `json:"organization"`
}

type IsRatingEnabled struct {
	IsEnabled bool          `json:"isEnabled"`
	Config    *RatingConfig `json:"config,omitempty"`
}

func RatingConfigFromModel(c *models.RatingOrgConfig) *RatingConfig {
	return &RatingConfig{
		StartHour: c.StartTime,
		Weekday:   uint8(c.WeekDay),
		Org:       OrgViewFromModelShort(c.Organization),
	}
}

type RatingQuestions struct {
	Initiative []*RatingQuestion `json:"initiative"`
	Discipline []*RatingQuestion `json:"discipline"`
	Efficiency []*RatingQuestion `json:"efficiency"`
	Teamwork   []*RatingQuestion `json:"teamwork"`
	Loyalty    []*RatingQuestion `json:"loyalty"`
}

type RatingQuestion struct {
	Text     string                  `json:"text"`
	Required bool                    `json:"required"`
	Options  []*RatingQuestionOption `json:"options,omitempty"`
	Range    *RatingResponseRange    `json:"range,omitempty"`
}

type RatingQuestionOption struct {
	Text string `json:"text"`
	Cost int    `json:"cost"`
}

type RatingResponseRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

func RatingEventFromModel(model *models.RatingEvent) *Rating {
	return &Rating{
		ID:           model.ID,
		Start:        model.Start.Unix(),
		End:          model.End.Unix(),
		Organization: OrgViewFromModelShort(&model.Organization),
	}
}
