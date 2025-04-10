package controller

import (
	"context"

	"github.com/Fl0rencess720/Springbroad/internal/data"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PortfolioRepo interface {
	GetAllTemplatesFromDB(context.Context) ([]data.Template, error)
	GetAllTemplatesFromRedis(context.Context) ([]data.Template, error)
	SaveAllTemplatesToRedis(context.Context, []data.Template) error
	GetHotTemplatesFromDB(context.Context) ([]data.Template, error)
	GetHotTemplatesFromRedis(context.Context) ([]data.Template, error)
	GetPortfolioFromDB(context.Context, string) ([]data.Portfolio, error)
	GetPortfolioFromRedis(context.Context, string) ([]data.Portfolio, error)
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

}

func (uc *PortfolioUsecase) SavePortfolio(c *gin.Context) {

}
