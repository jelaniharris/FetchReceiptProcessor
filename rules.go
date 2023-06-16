package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func alphanumericPoints(str string) int {
	// Only consider NOT alphanumeric characters
	var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

	// In the string replace the matching ones with a blank value
	newStr := nonAlphanumericRegex.ReplaceAllString(str, "")

	return len(newStr)
}

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

func itemsLengthGrouping(items []item, grouping int) int {
	numberOfGroups := len(items) / grouping
	return numberOfGroups
}

func itemDescriptionPricePoints(item item) (int, float64) {
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

func oddPurchaseDate(str string) bool {
	// Magical reference time
	tt, error := time.Parse("2006-01-02", str)

	if error != nil {
		return false
	}

	if tt.Day()%2 != 0 {
		return true
	}

	return false
}

func checkPurchaseTime(str string) bool {
	var hour int
	var minute int
	n, err := fmt.Sscanf(str, "%d:%d", &hour, &minute)

	if err != nil {
		return false
	}

	if n > 0 && hour >= 14 && hour < 16 {
		return true
	}

	return false
}
