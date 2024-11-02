package main

import (
	"fmt"
	"net/http"
)

// A map to store receipts in memory.
// Based on requirements we could store the points only.
// I decided to store the receipts as well, thinking that for future project enhancements it will be useful.
var receipts = make(map[string]ReceiptEnhanced)

func main() {

	// Handlers for processing receipts and retrieving points for a specific receipt
	http.HandleFunc("/receipts/process", processReceiptHandler)
	http.HandleFunc("/receipts/", getPointsHandler)

	// Start the HTTP server on port 8080
	fmt.Println("Starting server on :8080...")
	http.ListenAndServe(":8080", nil)
}
