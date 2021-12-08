package model

import (
	"errors"

	"gorm.io/gorm"
)

type Question struct {
	gorm.Model
	QuestionID string `gorm:"Index"`
	Question   string
	Answer     string
	Readness   uint8
}

func ListQuestions() ([]Question, error) {
	qlist := []Question{}
	query := GetInstance().db.Model(&Question{}).Where("readness=0")
	err := query.Find(&qlist).Error
	if err != nil {
		return nil, err
	}
	return qlist, nil
}

func GetQuestion(questionId string) (*Question, error) {
	q := Question{}
	err := GetInstance().db.Model(&Question{}).Where("question_id=? AND readness=0", questionId).Find(&q).Error
	if err != nil {
		return nil, err
	}
	if q.ID == 0 {
		return nil, errors.New("this Question Has been READ")
	}
	return &q, nil
}

func CreateQuestions(quesList []Question) error {
	return GetInstance().db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(quesList).Error
	})
}

func UpdateQuestionReadness(questionId string) error {
	return GetInstance().db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&Question{}).Where("question_id=?", questionId).Update("readness", 1).Error
	})
}
