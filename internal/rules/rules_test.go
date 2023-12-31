package rules

import (
	"errors"
	"testing"

	"github.com/jelaniharris/FetchReceiptProcessor/internal/models"
)

// Structures for testing point generation rules
type AlphanumericLengthStruct struct {
	arg1     string
	expected int
}

type IsTotalRoundStruct struct {
	arg1     string
	expected bool
}

type IsTotalAMultiplierStruct struct {
	arg1     string
	arg2     float64
	expected bool
}

type ItemsLengthGroupingStruct struct {
	arg1     []models.Item
	arg2     int
	expected int
}

type ItemDescriptionPricePointsStruct struct {
	arg1           models.Item
	expectedLength int
	expectedValue  float64
}

type OddPurchaseDateStruct struct {
	arg1     string
	expected bool
}

type CheckPurchaseTimeStruct struct {
	arg1     string
	expected bool
	errorMsg error
}

func TestAlphanumericLength(t *testing.T) {

	testTable := []AlphanumericLengthStruct{
		{"abcd", 4},
		{"abcdefghijklmnopqrstuvwxyz", 26},
		{"ABDEFGHIJK", 10},
		{"0123456789", 10},
		{"True Value", 9},
		{"Target", 6},
		{"Kohl's", 5},
		{"J.C. Penney", 8},
		{"M&M Corner Market", 14},
	}

	for _, test := range testTable {
		if output := alphanumericLength(test.arg1); output != test.expected {
			t.Errorf("alphanumericLength(%q) = got %d, wanted %d", test.arg1, output, test.expected)
		}
	}
}

func TestIsTotalRound(t *testing.T) {
	testTable := []IsTotalRoundStruct{
		{"10.00", true},
		{"5.45", false},
		{"6.01", false},
		{"3400.00", true},
	}

	for _, test := range testTable {
		if output := isTotalRound(test.arg1); output != test.expected {
			t.Errorf("isTotalRound(%q) = got %t, wanted %t", test.arg1, output, test.expected)
		}
	}
}

func TestIsTotalAMultiplier(t *testing.T) {
	testTable := []IsTotalAMultiplierStruct{
		{"10.00", 0.25, true},
		{"5.04", 0.25, false},
		{"2.75", 0.25, true},
		{"7400.00", 0.25, true},
		{"9.01", 0.25, false},
	}

	for _, test := range testTable {
		if output := isTotalAMultiplier(test.arg1, test.arg2); output != test.expected {
			t.Errorf("isTotalAMultiplier(%q, %f) = got %t, wanted %t", test.arg1, test.arg2, output, test.expected)
		}
	}
}

func TestItemsLengthGrouping(t *testing.T) {
	testTable := []ItemsLengthGroupingStruct{
		{nil, 2, 0},
		{[]models.Item{}, 2, 0},
		{[]models.Item{
			{ShortDescription: "Pepsi - 12-oz", Price: "1.25"},
			{ShortDescription: "Dasani", Price: "1.40"},
		}, 2, 1},
		{[]models.Item{
			{ShortDescription: "Pepsi - 12-oz", Price: "1.25"},
		}, 2, 0},
		{[]models.Item{
			{ShortDescription: "Pepsi - 12-oz", Price: "1.25"},
			{ShortDescription: "Dasani", Price: "1.40"},
			{ShortDescription: "Mike & Ikes", Price: "1.15"},
			{ShortDescription: "Snickers Ice Cream Bar", Price: "2.25"},
		}, 2, 2},
		{[]models.Item{
			{ShortDescription: "Vitamin Water", Price: "1.99"},
			{ShortDescription: "Mike & Ikes", Price: "1.15"},
			{ShortDescription: "Snickers Ice Cream Bar", Price: "2.25"},
		}, 2, 1},
	}

	for _, test := range testTable {
		if output := itemsLengthGrouping(test.arg1, test.arg2); output != test.expected {
			t.Errorf("itemsLengthGrouping(%q, %d) = got %d, wanted %d", test.arg1, test.arg2, output, test.expected)
		}
	}
}

func TestItemDescriptionPricePoints(t *testing.T) {
	testTable := []ItemDescriptionPricePointsStruct{
		{models.Item{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}, 13, 0},
		{models.Item{ShortDescription: "Target", Price: "1.25"}, 6, 0.25},
		{models.Item{ShortDescription: "  Pez  ", Price: "1.00"}, 3, 0.20},
	}

	for _, test := range testTable {
		outputLength, outputValue := itemDescriptionPricePoints(test.arg1)

		if outputLength != test.expectedLength {
			t.Errorf("itemDescriptionPricePoints(%q) = got %d length, wanted %d length", test.arg1.ShortDescription, outputLength, test.expectedLength)
		}

		if outputValue != test.expectedValue {
			t.Errorf("itemDescriptionPricePoints(%q) = got %f value, wanted %f value", test.arg1.ShortDescription, outputValue, test.expectedValue)
		}
	}
}

func TestOddPurchaseDate(t *testing.T) {
	testTable := []OddPurchaseDateStruct{
		{"2022-06-01", true},
		{"2022-06-15", true},
		{"2022-06-16", false},
		{"2022-06-08", false},
	}

	for _, test := range testTable {
		if output := oddPurchaseDate(test.arg1); output != test.expected {
			t.Errorf("TestOddPurchaseDate(%q) = got %t, wanted %t", test.arg1, output, test.expected)
		}
	}
}

func TestCheckPurchaseTime(t *testing.T) {
	testTable := []CheckPurchaseTimeStruct{
		{"11:15", false, nil},
		{"14:00", false, nil},
		{"14:01", true, nil},
		{"15:59", true, nil},
		{"14:59", true, nil},
		{"15:00", true, nil},
		{"16:00", false, nil},
		{"23:59", false, nil},
		{"01:59", false, nil},
		{"00:00", false, nil},
		{"26:00", false, errors.New("Invalid Hour format")},
		{"23:64", false, errors.New("Invalid Minute format")},
	}

	for _, test := range testTable {
		output, err := checkPurchaseTime(test.arg1)
		if output != test.expected {
			t.Errorf("checkPurchaseTime(%q) = got %t, wanted %t", test.arg1, output, test.expected)
		}
		if err != nil && err.Error() != test.errorMsg.Error() {
			t.Errorf("checkPurchaseTime(%q) = expected error was %q, wanted %q", test.arg1, err.Error(), test.errorMsg)
		}
	}
}
