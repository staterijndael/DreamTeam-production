package views

import (
	"dt/models"
	"io/ioutil"
)

type FileMetaInfo struct {
	ID       uint   `json:"id"`
	Checksum string `json:"checksum"`
	Size     uint   `json:"size"`
}

type FileContent struct {
	Content string `json:"content"`
}

type File struct {
	*FileMetaInfo
	*FileContent
}

func FileMetaInfoViewFromModel(f *models.File) *FileMetaInfo {
	return &FileMetaInfo{
		ID:       f.ID,
		Checksum: f.Checksum,
		Size:     f.Size,
	}
}

func FileContentViewFromModel(f *models.File) (*FileContent, error) {
	buf, err := ioutil.ReadFile(f.FilePath)
	if err != nil {
		return nil, err
	}

	return &FileContent{
		Content: string(buf),
	}, nil
}

func FileViewFromModel(f *models.File) (*File, error) {
	content, err := FileContentViewFromModel(f)
	if err != nil {
		return nil, err
	}

	return &File{
		FileMetaInfo: FileMetaInfoViewFromModel(f),
		FileContent:  content,
	}, nil
}
