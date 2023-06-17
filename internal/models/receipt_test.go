package models

import (
	"testing"
)

func TestEmptyGetReceiptById(t *testing.T) {
	t.Cleanup(resetState)

	_, error := GetReceiptById("555")

	if error == nil {
		t.Errorf("GetReceiptById should throw an error if id doesn't exist")
	}
}

func TestAddReceiptsAndGetReceiptsById(t *testing.T) {
	t.Cleanup(resetState)

	newReceipt := Receipt{
		Retailer:     "Target",
		PurchaseDate: "2023-06-16",
		PurchaseTime: "13:30",
		Total:        "0.00",
		Items:        nil,
	}

	// Add a new receipt to our list
	newId := AddToReceipts(newReceipt)

	// Retrieve that receipt using the ID
	foundReceipt, error := GetReceiptById(newId)

	if error != nil {
		t.Errorf("GetReceiptById should have found the id")
	}

	if foundReceipt.ID != newId {
		t.Errorf("GetReceiptById did not find the id it just added")
	}
}

func TestGetReceipts(t *testing.T) {
	t.Cleanup(resetState)

	emptyReceipts := GetReceipts()

	if len(emptyReceipts) != 0 {
		t.Errorf("GetReceipts should be empty")
	}

	newReceipt := Receipt{
		Retailer:     "Target",
		PurchaseDate: "2023-06-16",
		PurchaseTime: "13:30",
		Total:        "0.00",
		Items:        nil,
	}

	// Add several new receipts to our list
	for n := 0; n < 5; n++ {
		AddToReceipts(newReceipt)
	}

	allReceipts := GetReceipts()
	receiptCount := len(allReceipts)

	if receiptCount != 5 {
		t.Errorf("GetReceipts did not get the corrent number of receipts: expected %d, got %d", 5, receiptCount)
	}

}

func resetState() {
	ClearReceipts()
}
