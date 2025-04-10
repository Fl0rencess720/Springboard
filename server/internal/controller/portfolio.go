package controller

import (
	"context"

	"github.com/Fl0rencess720/Springbroad/internal/data"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SavePortfolioRequest struct {
	Portfolio data.Portfolio `json:"portfolio"`
	Works     []data.Work    `json:"works"`
}

type PortfolioRepo interface {
	GetAllTemplatesFromDB(context.Context) ([]data.Template, error)
	GetAllTemplatesFromRedis(context.Context) ([]data.Template, error)
	SaveAllTemplatesToRedis(context.Context, []data.Template) error
	GetHotTemplatesFromDB(context.Context) ([]data.Template, error)
	GetHotTemplatesFromRedis(context.Context) ([]data.Template, error)
	IncreTemplateScore(context.Context, string) error
	GetPortfolioFromDB(context.Context, string) ([]data.Portfolio, error)
	GetPortfolioFromRedis(context.Context, string) ([]data.Portfolio, error)
	SavePortfolioToDB(context.Context, data.Portfolio, []data.Work) error
}

type PortfolioUsecase struct {
	repo PortfolioRepo
}

func NewPortfolioUsecase(repo PortfolioRepo) *PortfolioUsecase {
	return &PortfolioUsecase{repo: repo}
}

func (uc *PortfolioUsecase) GetAllTemplates(c *gin.Context) {
	templates, err := uc.repo.GetAllTemplatesFromRedis(c)
	if err == nil {
		SuccessResponse(c, templates)
		return
	}
	zap.L().Error("GetAllTemplatesFromRedis error", zap.Error(err))

	templates, err = uc.repo.GetAllTemplatesFromDB(c)
	if err != nil {
		ErrorResponse(c, ServerError, err)
		return
	}
	if err := uc.repo.SaveAllTemplatesToRedis(c, templates); err != nil {
		zap.L().Error("SaveAllTemplatesToRedis error", zap.Error(err))
	}
	SuccessResponse(c, templates)
}

func (uc *PortfolioUsecase) GetHotTemplates(c *gin.Context) {
	templates, err := uc.repo.GetHotTemplatesFromRedis(c)
	if err == nil {
		SuccessResponse(c, templates)
		return
	}
	zap.L().Error("GetHotTemplatesFromRedis error", zap.Error(err))
	templates, err = uc.repo.GetHotTemplatesFromDB(c)
	if err != nil {
		ErrorResponse(c, ServerError, err)
		return
	}
	SuccessResponse(c, templates)
}

func (uc *PortfolioUsecase) SavePortfolio(c *gin.Context) {
	req := SavePortfolioRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, ServerError, err)
		return
	}
	for i := 0; i < len(req.Works); i++ {
		req.Works[i].OSSKey = uuid.New().String()
	}
	if err := uc.repo.SavePortfolioToDB(c, req.Portfolio, req.Works); err != nil {
		ErrorResponse(c, ServerError, err)
		return
	}
	if err := uc.repo.IncreTemplateScore(c, req.Portfolio.TemplateUID); err != nil {
		zap.L().Error("IncreTemplateScore error", zap.Error(err))
	}
	SuccessResponse(c, gin.H{
		"works": req.Works,
	})
}
