package routes

import (
	"github.com/gin-gonic/gin"
	controller "golang-restrogo/controllers"
)

func OrderRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/orders", controller.GetOrders())
	incomingRoutes.GET("/orders/:order_id", controller.GetOrder())
	incomingRoutes.POST("/orders", controller.CreateOrder())
	incomingRoutes.POST("/orders/:order_id", controller.UpdateOrder())
}
