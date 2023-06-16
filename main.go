package main

import (
	"github.com/jelaniharris/FetchReceiptProcessor/internal/api"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	receiptsGroup := router.Group("/receipts")

	// Get a listing of all receipts
	receiptsGroup.GET("", api.GetReceipts)
	// Get a single receipt by an id
	receiptsGroup.GET(":id", api.GetReceipt)
	// Return the point value of a receipt
	receiptsGroup.GET(":id/points", api.GetReceiptPoints)
	// Creates a receipt
	receiptsGroup.POST("process", api.CreateReceipt)

	// Start up the server at 8080
	router.Run("localhost:8080")
}
