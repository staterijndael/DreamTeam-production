package rating

import (
	"dt/config"
	"dt/events"
	"dt/logwrap"
	"dt/managers/eventEmitter"
	"dt/models"
	"github.com/jinzhu/gorm"
	"time"
)

type Manager struct {
	db              *gorm.DB
	emitter         *eventEmitter.EventEmitter
	config          *config.Config
	day             time.Weekday
	hour            int
	listenersAmount uint32
}

func New(db *gorm.DB, emitter *eventEmitter.EventEmitter, config *config.Config) *Manager {
	rm := &Manager{
		db:      db,
		emitter: emitter,
		config:  config,
		day:     time.Now().Weekday(),
		hour:    time.Now().Hour(),
	}

	go rm.Loop()
	return rm
}

func (manager *Manager) findRatings(hour int, day time.Weekday) []*models.RatingOrgConfig {
	var ratingConfigs []*models.RatingOrgConfig
	if err := manager.db.
		Where("week_day = ? and start_time = ?", day, hour).
		Find(&ratingConfigs).Error; err != nil {
		logwrap.Error("Manager.Loop() select RatingOrgConfig from DB error: %s", err.Error())
		return nil
	}

	return ratingConfigs
}

func (manager *Manager) timerListener(c <-chan time.Time) {
	for {
		<-c
		now := time.Now()
		configs := manager.findRatings(now.Hour(), now.Weekday())
		for _, conf := range configs {
			now = time.Now()
			event := &events.RatingStartedEvent{
				Start:          now,
				End:            now.Add(time.Duration(time.Hour.Nanoseconds() * int64(manager.config.RatingEventDuration))),
				OrganizationID: conf.OrganizationID,
			}

			manager.emitter.Emit(event)

		}
	}
}

func (manager Manager) Loop() {
	now := time.Now()
	time.Sleep(
		time.Until(
			time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location()),
		),
	)

	timer := time.NewTimer(time.Hour)
	for i := 0; i < 10; i++ {
		go manager.timerListener(timer.C)
	}
}
