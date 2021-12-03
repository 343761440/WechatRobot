package model

import (
	"time"

	"gorm.io/gorm"
)

type WxUser struct {
	gorm.Model
	UserId   string `gorm:"uniqueIndex;type:varchar(32)"`
	NickName string //备注名
	UserType int
	Birthday string
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

func ListWxUsers() ([]WxUser, error) {
	var res []WxUser
	query := GetInstance().db.Model(&WxUser{}).Where("user_type != ?", USER_ADMIN).Find(&res)
	if query.Error != nil {
		return nil, query.Error
	}
	return res, nil
}

func (wu *WxUser) IsBirthday() bool {
	if len(wu.Birthday) > 0 {
		if time.Now().Format("01-02") == wu.Birthday {
			return true
		}
	}
	return false
}

func (wu *WxUser) GetUserType() string {
	if wu.UserType == USER_NORMAL {
		return "Normal"
	} else if wu.UserType == USER_FRIEND {
		return "Friend"
	} else if wu.UserType == USER_IMPORTANT {
		return "Important"
	} else if wu.UserType == USER_ADMIN {
		return "Root"
	}
	return "Unknow"
}
