package api

import (
	"github.com/gin-gonic/gin"

	"github.com/IgorRamosBR/g73-techchallenge-order/internal/controllers"
)

type ApiParams struct {
	ProductController controllers.ProductController
	OrderController   controllers.OrderController
}

func NewApi(params ApiParams) *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.GET("/products", params.ProductController.GetProducts)
		v1.POST("/products", params.ProductController.CreateProducts)
		v1.PUT("/products/:id", params.ProductController.UpdateProduct)
		v1.DELETE("/products/:id", params.ProductController.DeleteProduct)

		v1.GET("/orders", params.OrderController.GetAllOrders)
		v1.POST("/orders", params.OrderController.CreateOrder)
		v1.GET("/orders/:id/status", params.OrderController.GetOrderStatus)
		v1.PUT("/orders/:id/status", params.OrderController.UpdateOrderStatus)
	}

	return router
}
