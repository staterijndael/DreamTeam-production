package rating

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"dt/utils"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/semrush/zenrpc"
	"time"
)

//оценить поль-ля в рамках данного события оценивания.
//zenrpc:efficiency список очков начисленных за соответствующий вопрос по теме "Эффективность"
//zenrpc:loyalty список очков начисленных за соответствующий вопрос по теме "Лояльность"
//zenrpc:professionalism список очков начисленных за соответствующий вопрос по теме "Профессионализм"
//zenrpc:discipline список очков начисленных за соответствующий вопрос по теме "Дисциплина"
//zenrpc:53 данное событие рейтинга не найдено.
//zenrpc:54 данный пользователь в данном событии вами уже оценен.
//zenrpc:55 данное событие недоступно для оценивания (уже закрыто или еще не начато).
//zenrpc:56 нельзя оценить человека с которым вы не состоите в одной группе.
//zenrpc:60 оценка эффективности вне диапазона [4, 40]
//zenrpc:61 оценка лояльности вне диапазона [-5, 8]
//zenrpc:6 оценка профессионализма вне диапазона [-1, 5]
//zenrpc:63 оценка дисциплины вне диапазона [-5, 5]
//zenrpc:64 не на все обязательные вопросы был дан ответ
//zenrpc:return в случае удачного выполнения возвращает сообщение "ok"
func (s *Service) Estimate(
	ctx context.Context,
	uid,
	ratingID uint,
	initiative,
	discipline,
	efficiency,
	teamwork,
	loyalty []int32,
) (*common.CodeAndMessage, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)

	var rating models.RatingEvent
	if err := s.db.First(&rating, ratingID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.RatingEventNotFound, err, nil) // 53
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	now := time.Now()
	if !rating.Start.Before(now) || !rating.End.After(now) {
		return nil, errors.New(errors.RatingEventNotAvailable, nil, nil) // 55
	}

	existingEstimate := models.Estimate{EstimatorID: me.ID, EstimatedID: uid, EventID: ratingID}
	if err := s.db.
		Where(&existingEstimate).
		First(&existingEstimate).Error; err == nil {
		return nil, errors.New(errors.ThisUserAlreadyEstimatedByU, nil, nil) // 54
	} else if !gorm.IsRecordNotFoundError(err) {
		return nil, errors.New(errors.Internal, err, nil)
	}

	requiredQuestionsCount, err := s.getRequiredQuestionsCount()
	if err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	if (len(initiative) + len(discipline) + len(efficiency) + len(teamwork) + len(loyalty)) < requiredQuestionsCount {
		return nil, errors.New(errors.NotAllRequiredQuestionsAnswered, nil, nil) // 64
	}

	initiativeScore := utils.IntSliceSum(initiative)
	disciplineScore := utils.IntSliceSum(discipline)
	efficiencyScore := utils.IntSliceSum(efficiency)
	teamworkScore := utils.IntSliceSum(teamwork)
	loyaltyScore := utils.IntSliceSum(loyalty)

	if initiativeScore < 5 || initiativeScore > 25 {
		return nil, errors.New(errors.InitiativeScoreOutOfRange, nil, nil)
	}

	if disciplineScore < 5 || disciplineScore > 25 {
		return nil, errors.New(errors.DisciplineScoreOutOfRange, nil, nil)
	}

	if efficiencyScore < 5 || efficiencyScore > 25 {
		return nil, errors.New(errors.EfficiencyScoreOutOfRange, nil, nil)
	}

	if teamworkScore < 5 || teamworkScore > 25 {
		return nil, errors.New(errors.TeamworkScoreOutOfRange, nil, nil)
	}

	if loyaltyScore < 5 || loyaltyScore > 25 {
		return nil, errors.New(errors.LoyaltyScoreOutOfRange, nil, nil)
	}

	initiativeScore = utils.ConvertToScale(initiativeScore, 5, 25, 5)
	disciplineScore = utils.ConvertToScale(disciplineScore, 5, 25, 5)
	efficiencyScore = utils.ConvertToScale(efficiencyScore, 5, 25, 5)
	teamworkScore = utils.ConvertToScale(teamworkScore, 5, 25, 5)
	loyaltyScore = utils.ConvertToScale(loyaltyScore, 5, 25, 5)

	estimate := &models.Score{
		Initiative: initiativeScore,
		Discipline: disciplineScore,
		Efficiency: efficiencyScore,
		Teamwork:   teamworkScore,
		Loyalty:    loyaltyScore,
	}

	var groups []*models.Group
	if err := s.db.
		Where(&models.Group{OrganizationID: rating.OrganizationID}).
		Find(&groups).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	var group *models.Group = nil
	for _, g := range groups {
		if g.Community.Contains(me.ID) && g.Community.Contains(uid) {
			group = g
			break
		}
	}

	if groups == nil || len(groups) == 0 || group == nil {
		return nil, errors.New(errors.CantEstimateUserNotInMutualGroup, nil, nil) // 56
	}

	bytes, _ := json.Marshal(estimate)
	estimateModel := models.Estimate{
		EstimatorID:    me.ID,
		EstimatedID:    uid,
		GroupID:        group.ID,
		OrganizationID: rating.OrganizationID,
		EventID:        rating.ID,
		Estimate:       postgres.Jsonb{RawMessage: bytes},
	}

	if err := s.db.Create(&estimateModel).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	return common.ResultOK, nil
}
