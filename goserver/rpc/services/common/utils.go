package common

import (
	"dt/config"
	"dt/models"
	"dt/utils"
	"github.com/jinzhu/gorm"
	"regexp"
)

const (
//ConnectionContextKey  = "currentConnection"
//UserContextKey        = "currentUser"
//UserQueryKey          = "me"
//HTTPRequestContextKey = "Request"
)

var (
	SplitBySpacesRegex = regexp.MustCompile(`\s`)
)

type FileInput struct {
	Content string `json:"content"`
}

func CreateDBFile(db *gorm.DB, c *config.Config, content []byte) (*models.File, error) {
	checksum := utils.BytesChecksum(content)
	f, fp, err := utils.NewMediaFile(c.MediaDir)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	_, err = f.Write(content)
	if err != nil {
		return nil, err
	}

	dbFile := models.File{
		Checksum: checksum,
		Size:     uint(len(content)),
		FilePath: *fp,
	}

	creation := db.Create(&dbFile)
	if creation.Error != nil {
		return nil, err
	}

	if err := db.First(&dbFile, dbFile.ID).Error; err != nil {
		return nil, err
	}

	return &dbFile, nil
}
