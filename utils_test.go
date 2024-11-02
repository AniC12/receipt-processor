package main

import (
	"testing"
)

func TestValidateReceipt(test *testing.T) {
	invalidPurchaseDates := []string{"yyyy-mm-dd", "01-01-2022", "2022/01/01", "20220101", "2022-13-01"}
	invalidPurchaseTimes := []string{"1301", "25:00", "13:60", "1:01"}
	invalidTotalAmounts := []string{"18.7a", "18.7", "18", "18.789", "abc"}
	invalidItemPrices := []string{"6.4", "6", "6.789", "xyz"}

	testCases := []struct {
		name        string
		receipt     Receipt
		expectError bool
	}{
		{
			name: "Valid receipt",
			receipt: Receipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items: []Item{
					{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
					{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
				},
				Total: "18.74",
			},
			expectError: false,
		},
	}

	// Add tests for invalid PurchaseDate formats
	for _, date := range invalidPurchaseDates {
		testCases = append(testCases, struct {
			name        string
			receipt     Receipt
			expectError bool
		}{
			name: "Invalid PurchaseDate format: " + date,
			receipt: Receipt{
				Retailer:     "Target",
				PurchaseDate: date,
				PurchaseTime: "13:01",
				Items: []Item{
					{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
				},
				Total: "6.49",
			},
			expectError: true,
		})
	}

	// Add tests for invalid PurchaseTime formats
	for _, time := range invalidPurchaseTimes {
		testCases = append(testCases, struct {
			name        string
			receipt     Receipt
			expectError bool
		}{
			name: "Invalid PurchaseTime format: " + time,
			receipt: Receipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-01",
				PurchaseTime: time,
				Items: []Item{
					{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
				},
				Total: "6.49",
			},
			expectError: true,
		})
	}

	// Add tests for invalid Total formats
	for _, total := range invalidTotalAmounts {
		testCases = append(testCases, struct {
			name        string
			receipt     Receipt
			expectError bool
		}{
			name: "Invalid Total format: " + total,
			receipt: Receipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items: []Item{
					{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
				},
				Total: total,
			},
			expectError: true,
		})
	}

	// Add tests for invalid Item Price formats
	for _, price := range invalidItemPrices {
		testCases = append(testCases, struct {
			name        string
			receipt     Receipt
			expectError bool
		}{
			name: "Invalid Item Price format: " + price,
			receipt: Receipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items: []Item{
					{ShortDescription: "Mountain Dew 12PK", Price: price},
				},
				Total: "6.40",
			},
			expectError: true,
		})
	}

	for _, tt := range testCases {
		test.Run(tt.name, func(t *testing.T) {
			err := validateReceipt(tt.receipt)
			if (err != nil) != tt.expectError {
				t.Errorf("validateReceipt() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestCalculatePoints(test *testing.T) {
	testCases := []struct {
		name     string
		receipt  Receipt
		expected int
	}{
		{
			name: "Valid receipt with various points",
			receipt: Receipt{
				Retailer:     "Target",     // 6 points
				PurchaseDate: "2022-01-01", // 6 points
				PurchaseTime: "13:01",
				Items: []Item{ // 10 points
					{ShortDescription: "Apple", Price: "3.00"},
					{ShortDescription: "Banana", Price: "10.01"}, // 3 points
					{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
					{ShortDescription: "Emils Cheese Pizza", Price: "12.25"}, // 3 points
				},
				Total: "18.74",
			},
			expected: 28,
		},
		{
			name: "Round dollar total",
			receipt: Receipt{
				Retailer:     "Walmart",    // 7 points
				PurchaseDate: "2022-01-03", // 6 points
				PurchaseTime: "15:30",      // 10 points
				Items: []Item{ // 5 points
					{ShortDescription: "Apple", Price: "3.00"},
					{ShortDescription: "Banana", Price: "1.00"}, // 1 point
				},
				Total: "4.00", // 50 + 25 points
			},
			expected: 104,
		},
		{
			name: "Multiple of 0.25",
			receipt: Receipt{
				Retailer:     "M&M Store", // 7 points
				PurchaseDate: "2022-01-02",
				PurchaseTime: "14:45", // 10 points
				Items: []Item{
					{ShortDescription: "Cookies", Price: "2.25"},
				},
				Total: "2.25", // 25 points
			},
			expected: 42,
		},
		{
			name: "Odd day purchase date",
			receipt: Receipt{
				Retailer:     "Best Buy",   // 7 points
				PurchaseDate: "2022-01-05", // 6 points
				PurchaseTime: "10:30",
				Items: []Item{
					{ShortDescription: "Laptop", Price: "999.99"}, // 200 points
				},
				Total: "999.99",
			},
			expected: 213,
		},
		{
			name: "Purchase during bonus time",
			receipt: Receipt{
				Retailer:     "Costco", // 6 points
				PurchaseDate: "2022-01-06",
				PurchaseTime: "15:15", // 10 points
				Items: []Item{
					{ShortDescription: "Water Bottle", Price: "5.00"}, // 1 point
				},
				Total: "5.00", // 50 + 25 points
			},
			expected: 92,
		},
		{
			name: "Multiple items with bonus points",
			receipt: Receipt{
				Retailer:     "Whole Foods 2", // 11 points
				PurchaseDate: "2022-01-07",    // 6 points
				PurchaseTime: "13:45",         // 25 points
				Items: []Item{ // 5 points
					{ShortDescription: "Organic Apple", Price: "1.50"},
					{ShortDescription: "Banana", Price: "0.75"}, // 1 point
				},
				Total: "2.25",
			},
			expected: 48,
		},
	}

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			actual := calculatePoints(testCase.receipt)
			if actual != testCase.expected {
				test.Errorf("calculatePoints() = %v, expected %v", actual, testCase.expected)
			}
		})
	}
}
