package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Song struct {
	gorm.Model
	Name       string
	PlayUrl    string
	Singer     string
	Uploader   string    //上传者
	UploadTime time.Time //上传时间
}

func CreateSong(s *Song) error {
	return GetInstance().db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(s).Error
	})
}

func DeleteSong(name string) error {
	return GetInstance().db.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&Song{Name: name}).Error
	})
}

func ListSongs(uploadName string, limit int) ([]Song, error) {
	var songs []Song
	query := GetInstance().db.Model(&Song{}).Order("upload_time DESC")
	if len(uploadName) > 0 {
		query = query.Where("uploader like=?", uploadName)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	query = query.Find(&songs)
	if query.Error != nil {
		return nil, query.Error
	}
	return songs, nil
}

func ListRootSongs(limit int) ([]Song, error) {
	var name string
	query := GetInstance().db.Model(&WxUser{}).Where("wx_users.user_type=?", USER_ADMIN).Pluck("nick_name", &name)
	if query.Error != nil {
		return nil, query.Error
	}
	if len(name) == 0 {
		return nil, errors.New("rootName is nil")
	}
	return ListSongs(name, limit)
}
