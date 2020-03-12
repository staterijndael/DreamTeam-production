package rating

import (
	"dt/views"
	"encoding/json"
	"io/ioutil"
)

func (s *Service) getRequiredQuestionsCount() (int, error) {
	fileQuestionsBytes, err := ioutil.ReadFile(s.conf.RatingFile)
	if err != nil {
		return 0, err
	}

	var rQuestions views.RatingQuestions
	if err := json.Unmarshal(fileQuestionsBytes, &rQuestions); err != nil {
		return 0, err
	}

	var requiredQuestionsCount int
	var questions []*views.RatingQuestion
	questions = append(rQuestions.Initiative, rQuestions.Discipline...)
	questions = append(questions, rQuestions.Efficiency...)
	questions = append(questions, rQuestions.Teamwork...)
	questions = append(questions, rQuestions.Loyalty...)
	for _, q := range questions {
		if q.Required {
			requiredQuestionsCount++
		}
	}

	return requiredQuestionsCount, nil
}
