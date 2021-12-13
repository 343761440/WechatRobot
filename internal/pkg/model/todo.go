package model

import (
	"time"

	"gorm.io/gorm"
)

type TodoItem struct {
	gorm.Model
	ItemInfo    string
	FinishTime  time.Time
	FinishState uint8
}

const (
	TODO_WAIT_FINISH uint8 = 0
	TODO_FINISH      uint8 = 1
	TODO_ALL         uint8 = 2
)

func ListTodoItems(state uint8) ([]TodoItem, error) {
	qlist := []TodoItem{}
	query := GetInstance().db.Model(&TodoItem{})
	if state != TODO_ALL {
		query = query.Where("finish_state=?", state)
	}
	err := query.Find(&qlist).Error
	if err != nil {
		return nil, err
	}
	return qlist, nil
}

func GetTodoItem(id int) (*TodoItem, error) {
	t := TodoItem{}
	err := GetInstance().db.Model(&TodoItem{}).Where("id=?", id).Find(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func CreateTodoItems(todolist ...TodoItem) error {
	return GetInstance().db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(todolist).Error
	})
}

func UpdateTodoFinish(id int, fin int) error {
	return GetInstance().db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&TodoItem{}).Where("id=?", id).Update("finish_state", fin).Error
	})
}
