package controller

import (
	"context"

	"github.com/Fl0rencess720/Springbroad/internal/data"
	"github.com/Fl0rencess720/Springbroad/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SavePortfolioRequest struct {
	UID         string      `json:"uid"`
	Title       string      `json:"title"`
	TemplateUID string      `json:"template_uid"`
	Works       []data.Work `json:"works"`
}

type GetAllTemplatesResponse struct {
	UID    string `json:"uid"`
	Name   string `json:"name"`
	OSSKey string `json:"oss_key"`
}

type PortfolioRepo interface {
	GetAllTemplatesFromDB(context.Context) ([]data.Template, error)
	GetAllTemplatesFromRedis(context.Context) ([]data.Template, error)
	GetTemplatesFromDB(context.Context, []string) ([]data.Template, error)
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
	response := []GetAllTemplatesResponse{}
	for _, template := range templates {
		response = append(response, GetAllTemplatesResponse{
			UID:    template.UID,
			Name:   template.Name,
			OSSKey: template.OSSKey,
		})
	}
	SuccessResponse(c, response)
}

func (uc *PortfolioUsecase) GetHotTemplates(c *gin.Context) {
	templates, err := uc.repo.GetHotTemplatesFromRedis(c)
	if err == nil {
		uids := []string{}
		for _, template := range templates {
			uids = append(uids, template.UID)
		}
		templatesWithMeta, err := uc.repo.GetTemplatesFromDB(context.Background(), uids)
		if err != nil {
			ErrorResponse(c, ServerError, err)
			return
		}
		SuccessResponse(c, templatesWithMeta)
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
	flag := false
	if req.UID == "" {
		req.UID = uuid.New().String()
		flag = true
	}
	for i := 0; i < len(req.Works); i++ {
		if req.Works[i].OSSKey == "" {
			req.Works[i].OSSKey = uuid.New().String()
			req.Works[i].PortfolioUID = req.UID
		}
	}
	if err := uc.repo.SavePortfolioToDB(c, data.Portfolio{UID: req.UID, Title: req.Title,
		TemplateUID: req.TemplateUID,
		Works:       req.Works, Openid: c.GetString("openid")}, req.Works); err != nil {
		ErrorResponse(c, ServerError, err)
		return
	}
	if flag {
		if err := uc.repo.IncreTemplateScore(c, req.TemplateUID); err != nil {
			zap.L().Error("IncreTemplateScore error", zap.Error(err))
		}
	}
	SuccessResponse(c, gin.H{
		"uid":   req.UID,
		"works": req.Works,
	})
}

func (uc *PortfolioUsecase) GetMyPortfolio(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	claims, _, err := middleware.ParseToken(tokenString)
	if err != nil {
		ErrorResponse(c, ServerError, nil)
		zap.L().Error("ParseToken error", zap.Error(err))
	}
	portfolio, err := uc.repo.GetPortfolioFromRedis(c, claims.Openid)
	if err == nil {
		SuccessResponse(c, portfolio)
		return
	}
	ErrorResponse(c, ServerError, nil)
	zap.L().Error("GetPortfolioFromRedis error", zap.Error(err))
	portfolio, err = uc.repo.GetPortfolioFromDB(c, claims.Openid)
	if err != nil {
		ErrorResponse(c, ServerError, err)
		return
	}
	SuccessResponse(c, portfolio)
}

func (uc *PortfolioUsecase) GetHistoricalUsageTemplates(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	claims, _, err := middleware.ParseToken(tokenString)
	if err != nil {
		ErrorResponse(c, ServerError, nil)
		zap.L().Error("ParseToken error", zap.Error(err))
	}
	templates := []data.Template{}
	portfolio, err := uc.repo.GetPortfolioFromRedis(c, claims.Openid)
	if err == nil {
		for _, i := range portfolio {
			templates = append(templates, i.Template)
		}
		SuccessResponse(c, templates)
		return
	}
	zap.L().Error("GetPortfolioFromRedis error", zap.Error(err))
	portfolio, err = uc.repo.GetPortfolioFromDB(c, claims.Openid)
	if err != nil {
		ErrorResponse(c, ServerError, err)
		return
	}
	for _, i := range portfolio {
		templates = append(templates, i.Template)
	}
	SuccessResponse(c, templates)
}
