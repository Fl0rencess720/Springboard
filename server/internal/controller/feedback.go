package controller

import (
	"context"
	"strconv"
	"time"

	"github.com/Fl0rencess720/Springbroad/internal/data"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UpdateStatusRequest struct {
	UID    string `json:"uid"`
	Status int    `json:"status"`
}

type AddFeedbackRequest struct {
	Content string `json:"content"`
}

type FeedbackRepo interface {
	AddFeedbackToDB(context.Context, data.Feedback) error
	GetAllFeedbacksFromDB(context.Context) ([]data.Feedback, error)
	GetFeedbacksByStatusFromDB(data.FeedbackStatus, context.Context) ([]data.Feedback, error)
	UpdateFeedbacksStatus(string, data.FeedbackStatus, context.Context) error
}

type FeedbackUseCase struct {
	repo FeedbackRepo
}

func NewFeedbackUseCase(repo FeedbackRepo) *FeedbackUseCase {
	return &FeedbackUseCase{
		repo: repo,
	}
}

func (sc *FeedbackUseCase) AddFeedback(c *gin.Context) {
	req := AddFeedbackRequest{}
	feedback := data.Feedback{}
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, ServerError, err)
		return
	}
	feedback.UID = uuid.New().String()
	feedback.Timestamp = time.Now()
	feedback.Content = req.Content
	if err := sc.repo.AddFeedbackToDB(c, feedback); err != nil {
		zap.L().Error("SaveFeedback error", zap.Error(err))
	}
	SuccessResponse(c, nil)
}

func (sc *FeedbackUseCase) GetAllFeedbacks(c *gin.Context) {
	feedbacks, err := sc.repo.GetAllFeedbacksFromDB(c)
	if err == nil {
		SuccessResponse(c, feedbacks)
		return
	}
	zap.L().Error("GetAllFeedbacksFromDB error", zap.Error(err))
}

func (sc *FeedbackUseCase) GetFeedbacksByStatus(c *gin.Context) {
	statusInt, err := strconv.Atoi(c.DefaultQuery("status", "0"))
	if err != nil {
		ErrorResponse(c, ServerError, err)
	}
	status := data.FeedbackStatus(statusInt)
	feedbacks, err := sc.repo.GetFeedbacksByStatusFromDB(status, c)
	if err == nil {
		SuccessResponse(c, feedbacks)
		return
	}
	zap.L().Error("GetFeedbacksByStatusFromDB error", zap.Error(err))
}

func (sc *FeedbackUseCase) UpdateFeedbacksStatus(c *gin.Context) {
	req := UpdateStatusRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, ServerError, err)
		return
	}
	status := data.FeedbackStatus(req.Status)
	if err := sc.repo.UpdateFeedbacksStatus(req.UID, status, c); err != nil {
		ErrorResponse(c, ServerError, err)
		return
	}
	SuccessResponse(c, nil)
}
