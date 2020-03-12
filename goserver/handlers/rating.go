package handlers

import (
	"dt/events"
	"dt/logwrap"
	"dt/models"
	"dt/scopes"
	"dt/utils"
	"dt/views"
	"encoding/json"
	"github.com/jinzhu/gorm/dialects/postgres"
	"io/ioutil"
	"time"
)

func ratingHandler(event interface{}) {
	var e *events.RatingStartedEvent
	var ok bool
	if e, ok = event.(*events.RatingStartedEvent); !ok {
		return
	}

	handleError := func(er error) {
		logwrap.Error(
			"ratingHandler. RatingStartedEvent{"+
				"Start: %d, End: %d, OrgID: %d}; error: %s",
			e.Start.Unix(),
			e.End.Unix(),
			e.OrganizationID,
			er.Error(),
		)
	}

	fileQuestionsBytes, err := ioutil.ReadFile(conf.RatingFile)
	if err != nil {
		handleError(err)
		return
	}

	ratingModel := models.RatingEvent{
		Start:          e.Start,
		End:            e.End,
		OrganizationID: e.OrganizationID,
		Questions:      postgres.Jsonb{RawMessage: fileQuestionsBytes},
	}

	if err := db.Create(&ratingModel).Error; err != nil {
		handleError(err)
		return
	}

	if err := db.First(&ratingModel, ratingModel.ID).Error; err != nil {
		handleError(err)
		return
	}

	var members []uint
	if err := db.
		Scopes(scopes.GroupMembersIDsOfOrg(ratingModel.OrganizationID, &members)).Error; err != nil {
		handleError(err)
		return
	}

	msg := &views.JSONRPCNotification{
		Method: "new",
		Params: &utils.Container{
			Type: "rating.started",
			Data: &struct {
				ID           uint       `json:"id"`
				Start        int64      `json:"start"`
				End          int64      `json:"end"`
				Organization *views.Org `json:"organization"`
			}{
				ID:           ratingModel.ID,
				Start:        e.Start.Unix(),
				End:          e.End.Unix(),
				Organization: views.OrgViewFromModelShort(&ratingModel.Organization),
			},
		},
	}

	plebeianSet := receiversToSet(members)
	for _, err := range sendToAllPlebeianMembers(receiversFromSet(plebeianSet), msg, nil) {
		handleError(err)
	}

	dashSet := receiversToSet(ratingModel.Organization.Admins.MembersIDs()).Union(plebeianSet)
	for _, err := range sendToAllOrgAdminMembers(receiversFromSet(dashSet), msg, nil) {
		handleError(err)
	}

	go processEstimates(ratingModel.End.Add(time.Second), ratingModel.ID)
	return
}

func processEstimates(end time.Time, ratingEventID uint) {
	time.Sleep(time.Until(end))

	handleError := func(msg string, er error) {
		logwrap.Error(
			"[handlers.rating.processEstimates] description: %s; error: %s",
			msg,
			er.Error(),
		)
	}

	var estimates []models.Estimate
	if err := db.
		Where("event = ?", ratingEventID).
		Find(&estimates).Error; err != nil {
		handleError("get estimates", err)
		return
	}

	//TODO можно ускорить, если не сохранять в базу каждый отдельный эстимейт, а сначала посчитать конечный скор
	// каждому оцененному, а потом всех засейвить
	usersToNotify := map[uint]models.User{}
	for _, estimate := range estimates {
		var score models.Score
		if err := json.Unmarshal(estimate.Estimate.RawMessage, &score); err != nil {
			handleError("unmarshal score", err)
			return
		}

		estimatedScore := estimate.Estimated.GetScore()
		estimatedScore.Adjust(&score)
		estimate.Estimated.SetScore(estimatedScore)
		if err := db.Model(&models.User{}).Where("id = ?", estimate.Estimated.ID).Update("score", estimate.Estimated.Score).Error; err != nil {
			handleError("save estimated to DB", err)
			return
		}

		usersToNotify[estimate.Estimated.ID] = estimate.Estimated
	}

	notifyUsers(usersToNotify, handleError)
}

func notifyUsers(users map[uint]models.User, handleError func(string, error)) {
	for id := range users {
		user := users[id]
		msg := &views.JSONRPCNotification{
			Method: "new",
			Params: &utils.Container{
				Type: "notification.userscorechanged",
				Data: &struct {
					Score views.Score `json:"score"`
				}{
					Score: views.ScoreViewFromModel(user.GetScore()),
				},
			},
		}

		for _, err := range sendToAllPlebeianMembers([]uint{users[id].ID}, msg, nil) {
			handleError("send notification", err)
		}
	}
}
