package controller

import (
	"github.com/Fl0rencess720/Springbroad/internal/data"
	"github.com/gin-gonic/gin"
)

type PortfolioRepo interface {
	GetAllTemplatesFromDB() []data.Template
	GetAllTemplatesFromRedis() []data.Template
	GetHotTemplatesFromDB() []data.Template
	GetHotTemplatesFromRedis() []data.Template
	GetPortfolioFromDB(openid string) []data.Portfolio
	GetPortfolioFromRedis(openid string) []data.Portfolio
}

type PortfolioUsecase struct {
	repo PortfolioRepo
}

func NewPortfolioUsecase(repo PortfolioRepo) *PortfolioUsecase {
	return &PortfolioUsecase{repo: repo}
}

func (uc *PortfolioUsecase) GetAllTemplates(c *gin.Context) {

}

func (uc *PortfolioUsecase) GetHotTemplates(c *gin.Context) {

}

func (uc *PortfolioUsecase) SavePortfolio(c *gin.Context) {

}
