package portfolio

import (
	"github.com/Fl0rencess720/Springbroad/internal/controller"
	"github.com/gin-gonic/gin"
)

func InitAPI(group *gin.RouterGroup, pu *controller.PortfolioUsecase) {
	group.GET("/template", pu.GetAllTemplates)
	group.GET("/template/hot", pu.GetHotTemplates)
	group.POST("/portfolio/save", pu.SavePortfolio)
	group.GET("/portfolio/me", pu.GetMyPortfolio)
	group.GET("/portfolio/history", pu.GetHistoricalUsageTemplates)
}
