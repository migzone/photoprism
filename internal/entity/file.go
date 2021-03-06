package entity

import (
	"fmt"
	"strings"

	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
	"github.com/photoprism/photoprism/internal/util"
)

// An image or sidecar file that belongs to a photo
type File struct {
	Model
	CreatedIn        int64
	UpdatedIn        int64
	Photo            *Photo
	PhotoID          uint   `gorm:"index;"`
	PhotoUUID        string `gorm:"index;"`
	FileUUID         string `gorm:"unique_index;"`
	FileName         string `gorm:"type:varchar(512);unique_index"` // max 3072 bytes / 4 bytes for utf8mb4 = 768 chars
	FileHash         string `gorm:"type:varchar(128);unique_index"`
	FileOriginalName string
	FileType         string `gorm:"type:varchar(32)"`
	FileMime         string `gorm:"type:varchar(64)"`
	FilePrimary      bool
	FileSidecar      bool
	FileVideo        bool
	FileMissing      bool
	FileDuplicate    bool
	FilePortrait     bool
	FileWidth        int
	FileHeight       int
	FileOrientation  int
	FileAspectRatio  float64
	FileMainColor    string `gorm:"type:varchar(32);index;"`
	FileColors       string
	FileLuminance    string
	FileChroma       uint
	FileNotes        string `gorm:"type:text"`
}

func FindFileByHash(db *gorm.DB, fileHash string) (File, error) {
	var file File

	q := db.Unscoped().First(&file, "file_hash = ?", fileHash)

	return file, q.Error
}

func (m *File) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("FileUUID", util.UUID())
}

func (m *File) DownloadFileName() string {
	if m.Photo == nil {
		return fmt.Sprintf("%s.%s", m.FileHash, m.FileType)
	}

	var name string

	if m.Photo.PhotoTitle != "" {
		name = strings.Title(slug.MakeLang(m.Photo.PhotoTitle, "en"))
	} else {
		name = m.PhotoUUID
	}

	taken := m.Photo.TakenAt.Format("20060102-150405")

	result := fmt.Sprintf("%s-%s.%s", taken, name, m.FileType)

	return result
}
