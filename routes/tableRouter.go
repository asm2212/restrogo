package routes

import (
	"github.com/gin-gonic/gin"
	controller "golang-restrogo/controllers"
)

func TableRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/tabels", controller.GetTables())
	incomingRoutes.GET("/tables/:table_id", controller.GetTable())
	incomingRoutes.POST("/tables", controller.CreateTable())
	incomingRoutes.POST("/tables/:table_id", controller.UpdateTable())
}
