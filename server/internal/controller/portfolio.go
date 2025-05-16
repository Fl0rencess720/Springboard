package controller

import (
	"context"

	"github.com/Fl0rencess720/Springboard/internal/data"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SavePortfolioRequest struct {
	UID         string         `json:"uid"`
	Title       string         `json:"title"`
	TemplateUID string         `json:"template_uid"`
	Projects    []data.Project `json:"projects"`
}

// type GetAllTemplatesResponse struct {
// 	UID    string `json:"uid"`
// 	Name   string `json:"name"`
// 	OSSKey string `json:"oss_key"`
// }

type PortfolioRepo interface {
	GetAllTemplatesFromDB(context.Context) ([]data.Template, error)
	GetAllTemplatesFromRedis(context.Context) ([]data.Template, error)
	GetTemplatesFromDB(context.Context, []string) ([]data.Template, error)
	SaveAllTemplatesToRedis(context.Context, []data.Template) error
	GetHotTemplatesFromDB(context.Context) ([]data.Template, error)
	GetHotTemplatesFromRedis(context.Context) ([]data.Template, error)
	IncreTemplateScore(context.Context, string) error
	GetTemplateByUIDFromDB(context.Context, string) (data.Template, error)

	GetPortfoliosFromDB(context.Context, string) ([]data.Portfolio, error)
	GetPortfoliosFromRedis(context.Context, string) ([]data.Portfolio, error)
	GetPortfolioByUIDFromDB(context.Context, string) (data.Portfolio, error)
	SavePortfoliosToRedis(context.Context, []data.Portfolio, string) error
	SavePortfolioToDB(context.Context, data.Portfolio) error
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
	// response := []GetAllTemplatesResponse{}
	// for _, template := range templates {
	// 	response = append(response, GetAllTemplatesResponse{
	// 		UID:    template.UID,
	// 		Name:   template.Name,
	// 		OSSKey: template.OSSKey,
	// 	})
	// }
	SuccessResponse(c, templates)
}

func (uc *PortfolioUsecase) GetTemplateByUID(c *gin.Context) {
	uid := c.Query("uid")
	template, err := uc.repo.GetTemplateByUIDFromDB(c, uid)
	if err != nil {
		ErrorResponse(c, ServerError, err)
		return
	}
	SuccessResponse(c, template)
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
	for i := 0; i < len(req.Projects); i++ {
		if req.Projects[i].UID == "" {
			req.Projects[i].UID = uuid.New().String()
			req.Projects[i].PortfolioUID = req.UID
		}
	}
	if err := uc.repo.SavePortfolioToDB(c, data.Portfolio{UID: req.UID, Title: req.Title,
		TemplateUID: req.TemplateUID,
		Projects:    req.Projects, Openid: c.GetString("openid")}); err != nil {
		ErrorResponse(c, ServerError, err)
		return
	}
	if flag {
		if err := uc.repo.IncreTemplateScore(c, req.TemplateUID); err != nil {
			zap.L().Error("IncreTemplateScore error", zap.Error(err))
		}
	}
	templates, err := uc.repo.GetTemplateByUIDFromDB(c, req.TemplateUID)
	if err != nil {
		zap.L().Error("GetTemplateByUIDFromDB error", zap.Error(err))
		ErrorResponse(c, ServerError, err)
		return
	}
	SuccessResponse(c, gin.H{
		"uid":      req.UID,
		"projects": req.Projects,
		"template": templates,
	})
}

func (uc *PortfolioUsecase) GetMyPortfolios(c *gin.Context) {
	openid := c.GetString("openid")
	portfolios, err := uc.repo.GetPortfoliosFromRedis(c, openid)
	if err == nil {
		SuccessResponse(c, portfolios)
		return
	}
	zap.L().Error("GetPortfolioFromRedis error", zap.Error(err))
	portfolios, err = uc.repo.GetPortfoliosFromDB(c, openid)
	if err != nil {
		ErrorResponse(c, ServerError, err)
		return
	}
	if err := uc.repo.SavePortfoliosToRedis(c, portfolios, openid); err != nil {
		zap.L().Error("SavePortfoliosToRedis error", zap.Error(err))
	}
	SuccessResponse(c, portfolios)
}

func (uc *PortfolioUsecase) GetPortfolioByUID(c *gin.Context) {
	uid := c.Query("uid")
	portfolio, err := uc.repo.GetPortfolioByUIDFromDB(c, uid)
	if err != nil {
		ErrorResponse(c, ServerError, err)
		return
	}
	SuccessResponse(c, portfolio)
}

func (uc *PortfolioUsecase) GetHistoricalUsageTemplates(c *gin.Context) {
	openid := c.GetString("openid")
	templates := []data.Template{}
	portfolios, err := uc.repo.GetPortfoliosFromRedis(c, openid)
	if err == nil {
		seen := make(map[string]struct{})
		for _, p := range portfolios {
			if _, ok := seen[p.TemplateUID]; !ok {
				templates = append(templates, p.Template)
				seen[p.TemplateUID] = struct{}{}
			}
		}
		SuccessResponse(c, templates)
		return
	}
	zap.L().Error("GetPortfolioFromRedis error", zap.Error(err))
	portfolios, err = uc.repo.GetPortfoliosFromDB(c, openid)
	if err != nil {
		ErrorResponse(c, ServerError, err)
		return
	}
	if err := uc.repo.SavePortfoliosToRedis(c, portfolios, openid); err != nil {
		zap.L().Error("SavePortfoliosToRedis error", zap.Error(err))
	}
	seen := make(map[string]struct{})
	for _, p := range portfolios {
		if _, ok := seen[p.TemplateUID]; !ok {
			templates = append(templates, p.Template)
			seen[p.TemplateUID] = struct{}{}
		}
	}
	SuccessResponse(c, templates)
}
