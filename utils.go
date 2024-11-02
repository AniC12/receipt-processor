package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// validateReceipt validates the receipt fields based on specific format requirements
func validateReceipt(receipt Receipt) error {
	// Validate PurchaseDate format to be YYYY-MM-DD
	layout := "2006-01-02"
	_, err := time.Parse(layout, receipt.PurchaseDate)
	if err != nil {
		return errors.New("invalid json. Please use YYYY-MM-DD format for PurchaseTime")
	}

	// Validate PurchaseTime format to be hh:mm
	re := regexp.MustCompile(`^([0-1][0-9]|2[0-3]):[0-5][0-9]$`)
	if !re.MatchString(receipt.PurchaseTime) {
		return errors.New("invalid json. Please use hh:mm format for PurchaseTime")
	}

	// Validate Total format to be a dollar amount with cents
	re = regexp.MustCompile(`^\d+(\.\d{2})$`)
	if !re.MatchString(receipt.Total) {
		return errors.New("invalid json. Total must represent a dollar amount with cents, like 120.40")
	}

	// Validate Price format to be a dollar amount with cents
	for _, item := range receipt.Items {
		if !re.MatchString(item.Price) {
			return errors.New("invalid json. Items price must represent a dollar amount with cents, like 120.40")
		}
	}

	return nil
}

// calculatePoints calculates the points for a given receipt based on specific rules
func calculatePoints(receipt Receipt) int {
	points := 0

	// One point for every alphanumeric character in the retailer name.
	for _, char := range receipt.Retailer {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			points++
		}
	}

	// 50 points if the total is a round dollar amount with no cents.
	if strings.HasSuffix(receipt.Total, ".00") {
		points += 50
	}

	// 25 points if the total is a multiple of 0.25.
	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err == nil && total != 0 && int(total*100)%25 == 0 {
		points += 25
	}

	// 5 points for every two items on the receipt.
	points += (len(receipt.Items) / 2) * 5

	// Points for items with description length multiple of 3.
	for _, item := range receipt.Items {
		trimmedDescription := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDescription)%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points += int(price*0.2 + 0.9999) // round up
		}
	}

	// 6 points if the day in the purchase date is odd.
	day := receipt.PurchaseDate[8:]
	if dayInt, _ := strconv.Atoi(day); dayInt%2 != 0 {
		points += 6
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	hour, _ := strconv.Atoi(receipt.PurchaseTime[:2])
	minute, _ := strconv.Atoi(receipt.PurchaseTime[3:])
	if hour == 15 || (hour == 14 && minute > 0) {
		points += 10
	}

	return points
}
