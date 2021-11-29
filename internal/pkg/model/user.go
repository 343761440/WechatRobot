package model

import "gorm.io/gorm"

type WxUser struct {
	gorm.Model
	UserId   string `gorm:"uniqueIndex;type:varchar(32)"`
	NickName string //备注名
}

func CreateWxUser(userId string) error {
	return GetInstance().db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&WxUser{UserId: userId}).Error
	})
}

func DeleteWxUser(userId string) error {
	return GetInstance().db.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&WxUser{UserId: userId}).Error
	})
}

func GetWxUser(userId string) (*WxUser, error) {
	m := WxUser{}
	err := GetInstance().db.Model(&WxUser{}).Where("user_id=?", userId).Find(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}
