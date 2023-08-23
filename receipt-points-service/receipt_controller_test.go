package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockValidator struct {
	mock.Mock
}

func (m *MockValidator) Struct(interface{}) error {
	args := m.Called()
	return args.Error(0)
}

func TestProcessReceipt(t *testing.T) {
	recorder := httptest.NewRecorder()
	requestBody := []byte(`{"retailer": "Sample Retailer", "purchaseDate": "2023-08-21", "purchaseTime": "12:00", "total": "100.00", "items": [{"shortDescription": "Item 1", "price": "50.00"}]}`)
	request, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(requestBody))

	mockValidator := new(MockValidator)
	mockValidator.On("Struct").Return(nil)

	processReceipt(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var receiptId ReceiptId
	json.NewDecoder(recorder.Body).Decode(&receiptId)
	assert.NotEmpty(t, receiptId.ID)
}

func TestGetReceiptPoints(t *testing.T) {
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/receipts/some-receipt-id/points", nil)

	getReceiptPoints(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var points Points
	json.NewDecoder(recorder.Body).Decode(&points)
	assert.Equal(t, 0, points.Points)
}
