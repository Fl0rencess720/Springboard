package feedback

import (
	"github.com/Fl0rencess720/Springboard/internal/controller"
	"github.com/gin-gonic/gin"
)

func InitAPI(group *gin.RouterGroup, sc *controller.FeedbackUseCase) {
	group.POST("/add", sc.AddFeedback)
	group.GET("/all", sc.GetAllFeedbacks)
	group.GET("", sc.GetFeedbacksByStatus)
	group.POST("/update", sc.UpdateFeedbacksStatus)
}
