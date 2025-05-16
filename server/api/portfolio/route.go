package portfolio

import (
	"github.com/Fl0rencess720/Springboard/internal/controller"
	"github.com/gin-gonic/gin"
)

func InitAPI(group *gin.RouterGroup, pu *controller.PortfolioUsecase) {
	group.GET("/template/all", pu.GetAllTemplates)
	group.GET("/template/", pu.GetTemplateByUID)
	group.GET("/template/hot", pu.GetHotTemplates)
	group.POST("/portfolio/save", pu.SavePortfolio)
	group.GET("/portfolio/me", pu.GetMyPortfolios)
	group.GET("/portfolio/", pu.GetPortfolioByUID)
	group.GET("/portfolio/history", pu.GetHistoricalUsageTemplates)
}
