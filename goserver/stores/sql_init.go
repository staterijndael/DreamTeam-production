package stores

import (
	"dt/config"
	"dt/logwrap"
	"dt/models"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

func InitDB(c *config.Config) (*gorm.DB, func(), error) {
	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s", c.PostgresHost, c.PostgresPort,
			c.PostgresUserName, c.PostgresDBName, c.PostgresPassword, c.PostgresSSLMode),
	)

	if err != nil {
		return nil, nil, err
	}

	db.SetLogger(new(logger))
	db = db.Set("gorm:auto_preload", true)
	migrate(db)
	initModels(db)
	initIndices(db)
	DefaultAvatars, err = initDefAvatars(db, c)
	if err != nil {
		return nil, nil, err
	}

	cleanUp := func() {
		if err := db.Close(); err != nil {
			logwrap.Error("error closing gorm: %s", err.Error())
		}
	}

	if c.ZomboDBOn {
		err = zdbInit(db, c, &models.User{}, &models.Organization{}, &models.Group{}, &models.Nickname{})
		if err != nil {
			logwrap.Error("[zdb init]: %s", err.(*pq.Error).Code)
			if !IsNeedVacuumFullErr(err) {
				return nil, nil, err
			}

			if err = db.Exec("vacuum full").Error; err != nil {
				return nil, nil, err
			}
		}
	}

	return db, cleanUp, nil
}
