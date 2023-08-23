package main

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Receipt struct {
	Retailer     string  `json:"retailer" validate:"required"`
	PurchaseDate string  `json:"purchaseDate" validate:"required"`
	PurchaseTime string  `json:"purchaseTime" validate:"required"`
	Total        float32 `json:"total,string" validate:"required"`
	Items        []Item  `json:"items" validate:"required,dive"`
}

type Item struct {
	ShortDescription string  `json:"shortDescription"`
	Price            float32 `json:"price,string"`
}

type ReceiptId struct {
	ID string `json:"id"`
}

type Points struct {
	Points int `json:"points"`
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

	validate := validator.New()

	if err := validate.Struct(receipt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	receiptID := uuid.New().String()

	receiptsLock.Lock()
	receipts[receiptID] = receipt
	receiptsLock.Unlock()

	pointsLock.Lock()
	points[receiptID] = CalculateTotalPoints(receipt)
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

	json.NewEncoder(w).Encode(Points{Points: receiptPoints})
}
