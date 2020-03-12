package stores

import (
	"dt/config"
	"dt/models"
	"dt/utils"
	"github.com/jinzhu/gorm"
	"os"
)

type defaultAvatars struct {
	User  models.File
	Org   models.File
	Group models.File
}

var (
	DefaultAvatars *defaultAvatars
)

func createFileInDBIfNotExists(db *gorm.DB, path string) (*models.File, error) {
	f, err := getDBFile(path)
	if err != nil {
		return nil, err
	}

	err = db.Where(&f).First(&f).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}

		creation := db.Create(&f)
		if creation.Error != nil {
			return nil, err
		}

		creation.Scan(&f)
	}

	return f, nil
}

func getDBFile(path string) (*models.File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	st, err := f.Stat()
	if err != nil {
		return nil, err
	}

	checksumString, err := utils.FileCheckSum(f)
	if err != nil {
		return nil, err
	}

	return &models.File{
		Checksum: *checksumString,
		Size:     uint(st.Size()),
		FilePath: path,
	}, nil
}

func initDefAvatars(db *gorm.DB, c *config.Config) (*defaultAvatars, error) {
	uv, err := createFileInDBIfNotExists(db, c.DefaultUserImagePath)
	if err != nil {
		return nil, err
	}

	ov, err := createFileInDBIfNotExists(db, c.DefaultOrgImagePath)
	if err != nil {
		return nil, err
	}

	gv, err := createFileInDBIfNotExists(db, c.DefaultGroupImagePath)
	if err != nil {
		return nil, err
	}

	return &defaultAvatars{
		User:  *uv,
		Org:   *ov,
		Group: *gv,
	}, nil
}
