package models

import "github.com/jinzhu/gorm"

type File struct {
	gorm.Model
	Checksum string `gorm:"column:checksum" `
	Size     uint   `gorm:"column:size"`
	FilePath string `gorm:"column:file_path"`
}

func (*File) TableName() string {
	return "files"
}
