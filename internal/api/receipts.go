package api

import (
	"net/http"

	"github.com/jelaniharris/FetchReceiptProcessor/internal/models"
	"github.com/jelaniharris/FetchReceiptProcessor/internal/rules"
	"github.com/jelaniharris/FetchReceiptProcessor/internal/types"

	"github.com/gin-gonic/gin"
)

// Get all of the receipts, no limit, no pagination
func GetReceipts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, models.GetReceipts())
}

// Get the receipt by id
func GetReceipt(c *gin.Context) {
	receipt, error := models.GetReceiptById(c.Param("id"))

	if error != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": error.Error()})
		return
	}

	c.JSON(http.StatusOK, receipt)
}

// Calculate the point value of the receipt
func GetReceiptPoints(c *gin.Context) {
	receipt, error := models.GetReceiptById(c.Param("id"))
	if error != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": error.Error()})
		return
	}

	points := rules.CalculatePoints(*receipt)

	c.JSON(http.StatusOK, gin.H{"points": points})
}

// Create a receipt and add it to our memory array of receipts
func CreateReceipt(c *gin.Context) {
	var newReceipt types.Receipt

	if err := c.BindJSON(&newReceipt); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error":   "VALIDATEERROR",
				"message": err.Error(),
			})

		return
	}

	// Add created receipt to list of receipts
	newId := models.AddToReceipts(newReceipt)

	c.JSON(http.StatusCreated, gin.H{"id": newId})
}
