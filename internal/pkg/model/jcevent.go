package model

import "gorm.io/gorm"

type JcEvent struct {
	gorm.Model
	Author string
	Event  string
}

func ListJCEvent() ([]JcEvent, error) {
	var res []JcEvent
	err := GetInstance().db.Model(&JcEvent{}).Find(&res).Error
	return res, err
}

func AddJCEvent(author, content string) error {
	return GetInstance().db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&JcEvent{
			Author: author,
			Event:  content,
		}).Error
	})
}
