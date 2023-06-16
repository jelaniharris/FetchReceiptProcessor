package models

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jelaniharris/FetchReceiptProcessor/internal/types"
)

// In-memory storage for the receipts
var receipts = []types.Receipt{}

// Searches the storage for a given receipt id and returns it
func GetReceiptById(id string) (*types.Receipt, error) {
	for i, b := range receipts {
		if b.ID == id {
			return &receipts[i], nil
		}
	}

	return nil, errors.New("Receipt not found")
}

// Empty the list of receipts
func ClearReceipts() {
	receipts = []types.Receipt{}
}

// Return our list of receipts
func GetReceipts() []types.Receipt {
	return receipts
}

// Add another receipt to our list of receipts
func AddToReceipts(newReceipt types.Receipt) string {
	newId := uuid.NewString()
	newReceipt.ID = newId
	receipts = append(receipts, newReceipt)
	return newId
}
