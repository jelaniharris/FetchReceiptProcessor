package api

import (
	"net/http"

	"github.com/jelaniharris/FetchReceiptProcessor/internal/models"
	"github.com/jelaniharris/FetchReceiptProcessor/internal/rules"

	"github.com/gin-gonic/gin"
)

// Error Message Info
// @Description Error Message Information
type ErrorMessage struct {
	// The message
	Message string `json:"message"`
}

// @Description Receipt points awarded response with points
type ReceiptPointsResponse struct {
	// The points awarded for the receipt
	Points int `json:"points" binding:"required"`
}

// @Description Receipt processed response with id
type CreatedReceiptResponse struct {
	// The new receipt id
	ID string `json:"id" binding:"required"`
}

// GetReceipts		godoc
// @Description 	Get all of the receipts, no limit, no pagination
// @Summary				Get All Receipts
// @Produce				application/json
// @Tags					reciepts
// @Success				200 {array} []models.Receipt{}
// @Router				/receipts [get]
func GetReceipts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, models.GetReceipts())
}

// GetReceipt			godoc
// @Description 	Get the receipt by id
// @Summary				Get A Receipt
// @Param					id path string true "get receipt by id"
// @Produce				application/json
// @Tags					reciepts
// @Success				200 {object} models.Receipt{} "success"
// @Failure				404 {object} ErrorMessage
// @Router				/receipts/{id} [get]
func GetReceipt(c *gin.Context) {
	receipt, error := models.GetReceiptById(c.Param("id"))

	if error != nil {
		c.IndentedJSON(http.StatusNotFound, ErrorMessage{Message: error.Error()})
		return
	}

	c.JSON(http.StatusOK, receipt)
}

// GetReceipt			godoc
// @Description 	Returns the points awarded for the receipt
// @Summary				Calculate Receipt Points
// @Param					id path string true "The ID of the receipt"
// @Produce				application/json
// @Tags					reciepts
// @Success				201 {object} ReceiptPointsResponse "The number of points awarded"
// @Failure				404 {object} ErrorMessage "No receipt found for that id"
// @Router				/receipts/{id}/points [get]
func GetReceiptPoints(c *gin.Context) {
	receipt, error := models.GetReceiptById(c.Param("id"))
	if error != nil {
		c.IndentedJSON(http.StatusNotFound, ErrorMessage{Message: error.Error()})
		return
	}

	points := rules.CalculatePoints(*receipt)

	c.JSON(http.StatusOK,
		ReceiptPointsResponse{Points: points})
}

// CreateReceipt	godoc
// @Description 	Create a receipt and add it to our memory array of receipts. Will not save the id if given one and will always make a new one.
// @Summary				Process Receipt
// @Param					receipt body models.Receipt true "new receipt to create"
// @Accept				application/json
// @Produce				application/json
// @Tags					reciepts
// @Success				201 {object} CreatedReceiptResponse "Returns the ID assigned to the receipt"
// @Failure				400 {object} ErrorMessage "The receipt is invalid"
// @Router				/receipts/process [post]
func CreateReceipt(c *gin.Context) {
	var newReceipt models.Receipt

	if err := c.BindJSON(&newReceipt); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ErrorMessage{Message: err.Error()})

		return
	}

	// Add created receipt to list of receipts
	newId := models.AddToReceipts(newReceipt)

	c.JSON(http.StatusCreated, CreatedReceiptResponse{ID: newId})
}
