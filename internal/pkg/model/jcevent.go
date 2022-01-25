package model

import "gorm.io/gorm"

type JcEvent struct {
	gorm.Model
	Author string
	Event  string
}

func ListJCEvent(author string) ([]JcEvent, error) {
	var res []JcEvent
	query := GetInstance().db.Model(&JcEvent{}).Debug()
	if len(author) > 0 {
		query = query.Where("author=?", author)
	}
	err := query.Find(&res).Error
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
