package main

// ReceiptEnhanced represents a receipt with additional information, including an ID and calculated points
type ReceiptEnhanced struct {
    Id  string
    Receipt      Receipt
    Points int
}

// Receipt represents the details of a purchase receipt
type Receipt struct {
    Retailer      string  `json:"retailer"`
    PurchaseDate  string  `json:"purchaseDate"`
    PurchaseTime  string  `json:"purchaseTime"`
    Items         []Item  `json:"items"`
    Total         string  `json:"total"`
}

// Item represents an item purchased in the receipt
type Item struct {
    ShortDescription string  `json:"shortDescription"`
    Price            string  `json:"price"`
}