package main

import (
	"github.com/jelaniharris/FetchReceiptProcessor/internal/api"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Get a listing of all receipts
	router.GET("/receipts", api.GetReceipts)
	// Get a single receipt by an id
	router.GET("/receipts/:id", api.GetReceipt)
	// Return the point value of a receipt
	router.GET("/receipts/:id/points", api.GetReceiptPoints)
	// Creates a receipt
	router.POST("/receipts/process", api.CreateReceipt)

	// Start up the server at 8080
	router.Run("localhost:8080")
}
