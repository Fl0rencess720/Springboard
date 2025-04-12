package data

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type FeedbackStatus int

const (
	Pending FeedbackStatus = iota
	Approved
	Rejected
)

type Feedback struct {
	UID       string         `json:"uid" gorm:"type:varchar(255)"`
	Content   string         `json:"content" gorm:"type:varchar(255)"`
	Timestamp time.Time      `json:"timestamp"`
	Status    FeedbackStatus `json:"status" gorm:"index;type:varchar(255)"`
}

type FeedbackRepo struct {
	mysqlDB *gorm.DB
}

func NewFeedbackRepo(mysqlDB *gorm.DB) *FeedbackRepo {
	return &FeedbackRepo{
		mysqlDB: mysqlDB,
	}
}

func (r *FeedbackRepo) AddFeedbackToDB(ctx context.Context, feedback Feedback) error {
	if err := r.mysqlDB.Create(&feedback).Error; err != nil {
		return err
	}
	return nil
}

func (r *FeedbackRepo) GetAllFeedbacksFromDB(ctx context.Context) ([]Feedback, error) {
	feedbacks := []Feedback{}
	if err := r.mysqlDB.Find(&feedbacks).Error; err != nil {
		return nil, err
	}
	return feedbacks, nil
}

func (r *FeedbackRepo) GetFeedbacksByStatusFromDB(status FeedbackStatus, ctx context.Context) ([]Feedback, error) {
	feedbacks := []Feedback{}
	if err := r.mysqlDB.Where("status = ?", status).Find(&feedbacks).Error; err != nil {
		return nil, err
	}
	return feedbacks, nil
}

func (r *FeedbackRepo) UpdateFeedbacksStatus(uid string, status FeedbackStatus, c context.Context) error {

	if err := r.mysqlDB.Model(&Feedback{}).Where("uid = ?", uid).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}
