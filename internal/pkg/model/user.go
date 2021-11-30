package model

import "gorm.io/gorm"

type WxUser struct {
	gorm.Model
	UserId   string `gorm:"uniqueIndex;type:varchar(32)"`
	NickName string //备注名
	UserType int
}

const (
	USER_NORMAL    int = 0 //普通用户
	USER_ADMIN     int = 1 //管理员
	USER_IMPORTANT int = 2 //重要的人
	USER_FRIEND    int = 3 //朋友
)

func CreateWxUserWithUserId(userId string) error {
	return GetInstance().db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&WxUser{UserId: userId}).Error
	})
}

func CreateWxUser(user *WxUser) error {
	return GetInstance().db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(user).Error
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
