package models

import (
	"errors"

	"github.com/google/uuid"
)

// An item is a purchased item on a receipt
type Item struct {
	// The Short Product Description for the item.
	ShortDescription string `json:"shortDescription" binding:"required"`
	// The total price payed for this item.
	Price string `json:"price" binding:"required"`
}

// A receipt is a listing of purchased items, and other metadata
type Receipt struct {
	// The ID of the receipt
	ID string `json:"id"`
	// The name of the retailer or store the receipt is from.
	Retailer string `json:"retailer" binding:"required"`
	// The date of the purchase printed on the receipt.
	PurchaseDate string `json:"purchaseDate" binding:"required" time_format:"2006-01-02"`
	// The time of the purchase printed on the receipt. 24-hour time expected.
	PurchaseTime string `json:"purchaseTime" binding:"required" time_format:"hh:mm"`
	// The total amount paid on the receipt.
	Total string `json:"total" binding:"required"`
	// The list of items in this receipt
	Items []Item `json:"items" binding:"required,dive"`
}

// In-memory storage for the receipts
var receipts = []Receipt{}

// Searches the storage for a given receipt id and returns it
func GetReceiptById(id string) (*Receipt, error) {
	for i, b := range receipts {
		if b.ID == id {
			return &receipts[i], nil
		}
	}

	return nil, errors.New("Receipt not found")
}

// Empty the list of receipts
func ClearReceipts() {
	receipts = []Receipt{}
}

// Return our list of receipts
func GetReceipts() []Receipt {
	return receipts
}

// Add another receipt to our list of receipts
func AddToReceipts(newReceipt Receipt) string {
	newId := uuid.NewString()
	newReceipt.ID = newId
	receipts = append(receipts, newReceipt)
	return newId
}
