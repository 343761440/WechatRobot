package model

import "gorm.io/gorm"

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
	err := GetInstance().db.Model(&Question{}).Where("question_id=?", questionId).Find(&q).Error
	if err != nil {
		return nil, err
	}
	return &q, nil
}

func CreateQuestions(quesList []Question) error {
	return GetInstance().db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(quesList).Error
	})
}
