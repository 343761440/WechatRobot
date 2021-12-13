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

func GetJCEvent(id uint) (JcEvent, error) {
	var res JcEvent
	err := GetInstance().db.Model(&JcEvent{}).Where("id=?", id).Limit(1).Find(&res).Error
	return res, err
}

func UpdateJCEvent(id uint, content string) error {
	return GetInstance().db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&JcEvent{}).Where("id=?", id).Update("event", content).Error
	})
}
