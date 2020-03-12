package rating

import (
	"dt/logwrap"
	"dt/rpc/services/errors"
	"dt/views"
	"encoding/json"
	"github.com/semrush/zenrpc"
	"io/ioutil"
)

func (s *Service) GetQuestions() (*views.RatingQuestions, *zenrpc.Error) {
	fileQuestionsBytes, err := ioutil.ReadFile(s.conf.RatingFile)
	if err != nil {
		logwrap.Error("[rating.GetQuestions] error while reading rating file: %s", err.Error())
		return nil, errors.New(errors.Internal, err, nil)
	}

	var result views.RatingQuestions
	err = json.Unmarshal(fileQuestionsBytes, &result)
	if err != nil {
		logwrap.Error("[rating.GetQuestions] error while unmarshal rating file: %s", err.Error())
		return nil, errors.New(errors.Internal, err, nil)
	}

	return &result, nil
}
