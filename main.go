package main

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type item struct {
	ShortDescription string `json:"shortDescription" binding:"required"`
	Price            string `json:"price" binding:"required"`
}

type receipt struct {
	ID           string `json:"id"`
	Retailer     string `json:"retailer" binding:"required"`
	PurchaseDate string `json:"purchaseDate" binding:"required" time_format:"2006-01-02"`
	PurchaseTime string `json:"purchaseTime" binding:"required" time_format:"hh:mm"`
	Total        string `json:"total" binding:"required"`
	Items        []item `json:"items" binding:"required"`
}

type PointRuleItem struct {
	Description       string
	Price             string
	DescriptionLength int
	Value             float64
}

type PointRules struct {
	AlphanumericPoints int
	RoundDollarPoints  int
	MultiplierPoints   int
	GroupingPoints     int
	GroupingAmount     int
	PurchaseDatePoints int
	PurchaseTimePoints int
	RuleItems          []PointRuleItem
}

var receipts = []receipt{}

func getReceipts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, receipts)
}

func getReceiptById(id string) (*receipt, error) {
	for i, b := range receipts {
		if b.ID == id {
			return &receipts[i], nil
		}
	}

	return nil, errors.New("Receipt not found")
}

func getReceipt(c *gin.Context) {
	receipt, error := getReceiptById(c.Param("id"))

	if error != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": error.Error()})
		return
	}

	c.JSON(http.StatusOK, receipt)
}

func calculatePoints(rec receipt) int {
	var rulePoints PointRules
	currentPoints := 0

	// One point for every alphanumeric character in the retailer name.
	rulePoints.AlphanumericPoints = alphanumericPoints(rec.Retailer)
	currentPoints += rulePoints.AlphanumericPoints

	// 50 points if the total is a round dollar amount with no cents.
	if isTotalRound(rec.Total) {
		rulePoints.RoundDollarPoints = 50
		currentPoints += rulePoints.RoundDollarPoints
	}

	// 5 points if the total is a multiple of 0.25
	if isTotalAMultiplier(rec.Total, 0.25) {
		rulePoints.MultiplierPoints = 25
		currentPoints += rulePoints.MultiplierPoints
	}

	// 5 points for every two items on the receipt.
	rulePoints.GroupingAmount = itemsLengthGrouping(rec.Items, 2)
	rulePoints.GroupingPoints = rulePoints.GroupingAmount * 5
	currentPoints += rulePoints.GroupingPoints

	// If the trimmed length of the item description is a multiple of 3,
	// multiply the price by 0.2 and round up to the nearest integer.
	// The result is the number of points earned.
	for _, item := range rec.Items {
		descrLength, value := itemDescriptionPricePoints(item)
		if value > 0 {
			rulePoints.RuleItems = append(rulePoints.RuleItems, PointRuleItem{Price: item.Price, Description: item.ShortDescription, DescriptionLength: descrLength, Value: value})
			currentPoints += int(math.Ceil(value))
		}
	}

	// 6 points if the day in the purchase date is odd.
	if oddPurchaseDate(rec.PurchaseDate) {
		rulePoints.PurchaseDatePoints = 6
		currentPoints += rulePoints.PurchaseDatePoints
	}

	if checkPurchaseTime(rec.PurchaseTime) {
		rulePoints.PurchaseTimePoints = 10
		currentPoints += rulePoints.PurchaseTimePoints
	}

	// Breakdown output
	showBreakdown(rulePoints, currentPoints, rec)

	return currentPoints
}

func showBreakdown(rulePoints PointRules, totalPoints int, rec receipt) {
	fmt.Println("Breakdown:")
	fmt.Printf("%6d points - Retailer name has %d alphanumeric characters \n", rulePoints.AlphanumericPoints, rulePoints.AlphanumericPoints)
	fmt.Printf("%6d points - Total is a round dollar amount \n", rulePoints.RoundDollarPoints)
	fmt.Printf("%6d points - Total is multiple of 0.25 \n", rulePoints.MultiplierPoints)
	fmt.Printf("%6d points - %d items (%d pairs @ %d points each) \n", rulePoints.GroupingPoints, len(rec.Items), rulePoints.GroupingAmount, 5)
	for _, ruleItem := range rulePoints.RuleItems {
		roundedPoints := int(math.Ceil(ruleItem.Value))
		fmt.Printf("%6d points - %q is %d characters (a multiple of 3)\n", roundedPoints, ruleItem.Description, ruleItem.DescriptionLength)
		fmt.Printf("%s item price of %q * 0.2 = %.2f, rounded up is %d points \n", strings.Repeat(" ", 16), ruleItem.Price, ruleItem.Value, roundedPoints)
	}
	fmt.Printf("%6d points - Purchase day is odd \n", rulePoints.PurchaseDatePoints)
	fmt.Printf("%6d points - %q is between 2:00pm and 4:00pm \n", rulePoints.PurchaseTimePoints, rec.PurchaseTime)
	fmt.Println("+ -------")
	fmt.Printf("= %d points", totalPoints)
}

func getReceiptPoints(c *gin.Context) {
	receipt, error := getReceiptById(c.Param("id"))
	if error != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": error.Error()})
		return
	}

	points := calculatePoints(*receipt)

	c.JSON(http.StatusOK, gin.H{"points": points})
}

func processReceipts(c *gin.Context) {
	var newReceipt receipt

	if err := c.BindJSON(&newReceipt); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error":   "VALIDATEERR-1",
				"message": err.Error(),
			})

		return
	}

	newReceipt.ID = uuid.NewString()

	receipts = append(receipts, newReceipt)
	c.JSON(http.StatusCreated, gin.H{"id": newReceipt.ID})
}

func getItems() {

}

func main() {
	router := gin.Default()

	router.GET("/receipts", getReceipts)
	router.GET("/receipts/:id", getReceipt)
	router.GET("/receipts/:id/points", getReceiptPoints)
	router.POST("/receipts/process", processReceipts)

	router.Run("localhost:8080")
}
