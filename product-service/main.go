package main

import (
	"github.com/Vaishnav-OM/product-service/db"
	"github.com/Vaishnav-OM/product-service/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	db.Init()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	handlers.InitCollection()
	r.POST("/products", handlers.CreateProduct)
	r.PUT("/products/:id", handlers.UpdateProduct)
	r.DELETE("/products/:id", handlers.DeleteProduct)
	r.GET("/products", handlers.GetProducts)
	r.GET("/products/:id", handlers.GetProductByID)
	//implement deleteall

	r.Run(":8080")
}
