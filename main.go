package main

import (
	"github.com/jelaniharris/FetchReceiptProcessor/internal/api"

	"github.com/gin-gonic/gin"
	_ "github.com/jelaniharris/FetchReceiptProcessor/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Fetch Receipt Processor API
// @description This is a webservice that allows you to create a receipt and calculate a point value of that receipt using a set of rules.
// @version 1.0

// @license.name  MIT

// @BasePath /
// @host localhost:8080

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	router := gin.Default()

	// Add swagger support
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	receiptsGroup := router.Group("/receipts")
	{
		// Get a listing of all receipts
		receiptsGroup.GET("", api.GetReceipts)
		// Get a single receipt by an id
		receiptsGroup.GET(":id", api.GetReceipt)
		// Return the point value of a receipt
		receiptsGroup.GET(":id/points", api.GetReceiptPoints)
		// Creates a receipt
		receiptsGroup.POST("process", api.CreateReceipt)
	}

	// Start up the server at 8080
	router.Run()
}
