package rules

import (
	"errors"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jelaniharris/FetchReceiptProcessor/internal/models"
)

// Stores the result of the point calculation of every item in a receipt
type PointRuleItem struct {
	Description       string
	Price             string
	DescriptionLength int
	Value             float64
}

// Stores the calculated results for a point calculation for a receipt
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

// Calculates the alphanumeric length of a string
func alphanumericLength(str string) int {
	// Only consider NOT alphanumeric characters
	var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

	// In the string replace the matching ones with a blank value
	newStr := nonAlphanumericRegex.ReplaceAllString(str, "")

	return len(newStr)
}

// Tests to see if the total in the string is a round number
// 10.00 = true, 10.25 = false
func isTotalRound(str string) bool {
	amount, error := strconv.ParseFloat(str, 64)

	if error != nil {
		return false
	}

	// Take the amount minus the double casted amount int -> float and if there's any left
	// Then it's not a round dollar amount
	if amount-float64(int64(amount)) > 0.0 {
		return false
	}

	return true
}

// Checks to see if the total is a multiplier of the given multiplier
// 8.00 is a multiplier of 0.25 but not of 0.75
func isTotalAMultiplier(str string, multiplier float64) bool {
	amount, error := strconv.ParseFloat(str, 64)

	if error != nil {
		return false
	}

	// Can't use % here for floats - need to use the math.Mod function
	if math.Mod(amount, multiplier) == 0 {
		return true
	}

	return false
}

// Returns the number of groups that can be formed from a list of items
// Given 8 items and a grouping of 2, that's 4 pairs
// Given 9 items anda  grouping of 4, that's 2 ... quads?
func itemsLengthGrouping(items []models.Item, grouping int) int {

	// Prevent
	if grouping == 0 {
		return 0
	}

	numberOfGroups := len(items) / grouping
	return numberOfGroups
}

// Determines the trimmed length of an items description, and the point value
// based on if the length is divisible by three
func itemDescriptionPricePoints(item models.Item) (int, float64) {
	newStr := strings.Trim(item.ShortDescription, " ")
	length := len(newStr)

	if length%3 == 0 {
		amount, error := strconv.ParseFloat(item.Price, 64)

		if error != nil {
			return length, 0
		}

		value := amount * 0.2
		return length, value
	}
	return length, 0
}

// Checks to see if the date is odd
// 2023-05-05 = True, 2023-05-08 = False
func oddPurchaseDate(str string) bool {
	// Parse the time using the Magical reference time
	tt, error := time.Parse("2006-01-02", str)

	if error != nil {
		return false
	}

	// If day is odd
	if tt.Day()%2 != 0 {
		return true
	}

	return false
}

// Check to see if the purchase time is between 2pm and 4pm
// Input is assumed to be in 24 hour format e.g. 15:40
func checkPurchaseTime(str string) (bool, error) {
	var hour int
	var minute int

	// Parsing the string into a hh:mm format
	n, err := fmt.Sscanf(str, "%d:%d", &hour, &minute)

	// If the scanning produced an error
	if err != nil {
		return false, errors.New("Could not scan time")
	}

	// If the parsed time is somehow not a 24 hour clock
	if hour > 24 || hour < 0 {
		return false, errors.New("Invalid Hour format")
	}

	// If the parsed minutes is not in minutes
	if minute > 59 || minute < 0 {
		return false, errors.New("Invalid Minute format")
	}

	// If we parsed data and:
	// The hour is greater than 2pm
	// The hour is 2pm, and the minute is at least 1 or more
	// The hour is less than 6pm
	if n > 0 && (hour > 14 || hour == 14 && minute >= 1) && hour < 16 {
		return true, nil
	}

	return false, nil
}

// Given a receipt, calculate the amount of points it's worth based on
// a series of rules
func CalculatePoints(rec models.Receipt) (int, error) {
	var rulePoints PointRules
	currentPoints := 0

	// One point for every alphanumeric character in the retailer name.
	rulePoints.AlphanumericPoints = alphanumericLength(rec.Retailer)
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

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	checkedTime, err := checkPurchaseTime(rec.PurchaseTime)
	if err != nil {
		return 0, err
	}

	if checkedTime {
		rulePoints.PurchaseTimePoints = 10
		currentPoints += rulePoints.PurchaseTimePoints
	}

	// Breakdown output in console
	showBreakdown(rulePoints, currentPoints, rec)

	return currentPoints, nil
}

// Log the results of the point calculation
func showBreakdown(rulePoints PointRules, totalPoints int, rec models.Receipt) {
	log.Printf("Breakdown for Receipt ID (%q):", rec.ID)
	if rulePoints.AlphanumericPoints > 0 {
		log.Printf("%6d points - Retailer name has %d alphanumeric characters \n", rulePoints.AlphanumericPoints, rulePoints.AlphanumericPoints)
	}
	if rulePoints.RoundDollarPoints > 0 {
		log.Printf("%6d points - Total is a round dollar amount \n", rulePoints.RoundDollarPoints)
	}
	if rulePoints.MultiplierPoints > 0 {
		log.Printf("%6d points - Total is multiple of 0.25 \n", rulePoints.MultiplierPoints)
	}
	if rulePoints.GroupingPoints > 0 {
		log.Printf("%6d points - %d items (%d pairs @ %d points each) \n", rulePoints.GroupingPoints, len(rec.Items), rulePoints.GroupingAmount, 5)
	}
	// Then loop through our rule items
	for _, ruleItem := range rulePoints.RuleItems {
		roundedPoints := int(math.Ceil(ruleItem.Value))
		if roundedPoints > 0 {
			log.Printf("%6d points - %q is %d characters (a multiple of 3)\n", roundedPoints, ruleItem.Description, ruleItem.DescriptionLength)
			log.Printf("%s item price of %q * 0.2 = %.2f, rounded up is %d points \n", strings.Repeat(" ", 16), ruleItem.Price, ruleItem.Value, roundedPoints)
		}
	}
	if rulePoints.PurchaseDatePoints > 0 {
		log.Printf("%6d points - Purchase day is odd \n", rulePoints.PurchaseDatePoints)
	}
	if rulePoints.PurchaseTimePoints > 0 {
		log.Printf("%6d points - %q is between 2:00pm and 4:00pm \n", rulePoints.PurchaseTimePoints, rec.PurchaseTime)
	}
	log.Println("+ -------")
	log.Printf("= %d points", totalPoints)
}
