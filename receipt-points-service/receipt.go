package main

import (
	"encoding/json"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Receipt struct {
	Retailer     string  `json:"retailer"`
	PurchaseDate string  `json:"purchaseDate"`
	PurchaseTime string  `json:"purchaseTime"`
	Total        float32 `json:"total,string"`
	Items        []Item  `json:"items"`
}

type Item struct {
	ShortDescription string  `json:"shortDescription"`
	Price            float32 `json:"price,string"`
}

type ReceiptId struct {
	ID string `json:"id"`
}

var (
	receipts     = make(map[string]Receipt)
	points       = make(map[string]int)
	receiptsLock sync.RWMutex
	pointsLock   sync.RWMutex
)

func processReceipt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var receipt Receipt

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&receipt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	receiptID := uuid.New().String()

	receiptsLock.Lock()
	receipts[receiptID] = receipt
	receiptsLock.Unlock()

	pointsLock.Lock()
	points[receiptID] = calculatePoints(receipt)
	pointsLock.Unlock()

	json.NewEncoder(w).Encode(ReceiptId{ID: receiptID})
}

func getReceiptPoints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	receiptID := params["id"]

	pointsLock.RLock()
	receiptPoints := points[receiptID]
	pointsLock.RUnlock()

	json.NewEncoder(w).Encode(map[string]int{"points": receiptPoints})
}

var startTime2PM = time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC)
var endTime4PM = time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)

func calculatePoints(receipt Receipt) int {
	points := 0

	for _, char := range receipt.Retailer {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			points++
		}
	}

	totalDollars := int(receipt.Total)
	if float32(totalDollars) == receipt.Total {
		points += 50
	}

	if math.Mod(float64(receipt.Total), 0.25) == 0 {
		points += 25
	}

	points += (len(receipt.Items) / 2) * 5

	for _, item := range receipt.Items {
		trimmedLen := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLen%3 == 0 {
			itemPoints := int(math.Ceil(float64(item.Price) * 0.2))
			points += itemPoints
		}
	}

	purchaseDate, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err == nil && purchaseDate.Day()%2 != 0 {
		points += 6
	}

	purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)
	if err == nil {
		if purchaseTime.After(startTime2PM) && purchaseTime.Before(endTime4PM) {
			points += 10
		}
	}

	return points
}
