package model

import "gorm.io/gorm"

type LoveJoke struct {
	gorm.Model
	Content  string
	Readness int
}

func GetALoveJoke() (string, error) {
	joke := LoveJoke{}
	err := GetInstance().db.Transaction(func(tx *gorm.DB) error {
		err := tx.Debug().Model(&LoveJoke{}).Where("readness=0").Limit(1).Find(&joke).Error
		if err != nil {
			return err
		}
		err = tx.Model(&LoveJoke{}).Where("id=?", joke.ID).Update("readness", 1).Error
		return err
	})
	return joke.Content, err
}

func AddLoveJoke(content string) error {
	return GetInstance().db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&LoveJoke{Content: content}).Error
	})
}
