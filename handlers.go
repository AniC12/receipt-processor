package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// processReceiptHandler handles POST requests to process a receipt
func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = validateReceipt(receipt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// With current implementation if the same receipt is uploaded multiple times, a new id will be generated each time.
	// We could also keep a receipt->id map, but the requirements did not mention this, hence we skipped it.

	receiptEnhanced := ReceiptEnhanced{
		Id:      uuid.New().String(),
		Receipt: receipt,
		Points:  calculatePoints(receipt),
	}
	receipts[receiptEnhanced.Id] = receiptEnhanced

	response := map[string]string{"id": receiptEnhanced.Id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getPointsHandler handles GET requests to retrieve points for a specific receipt
func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/receipts/")
	id = strings.TrimSuffix(id, "/points")

	receiptEnhanced, exists := receipts[id]
	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	points := receiptEnhanced.Points

	response := map[string]int{"points": points}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
