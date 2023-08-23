package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalculateTotalPoints(t *testing.T) {

	receipt := Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []Item{
			{
				ShortDescription: "Mountain Dew 12PK",
				Price:            6.49,
			},
			{
				ShortDescription: "Emils Cheese Pizza",
				Price:            12.25,
			},
			{
				ShortDescription: "Knorr Creamy Chicken",
				Price:            1.26,
			},
			{
				ShortDescription: "Doritos Nacho Cheese",
				Price:            3.35,
			},
			{
				ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
				Price:            12.00,
			},
		},
		Total: 35.35,
	}

	points := CalculateTotalPoints(receipt)

	assert.Equal(t, 28, points)
}

func TestCalculatePointsIsAlphanumeric(t *testing.T) {
	points := calculatePointsIsAlphanumeric("Sample@Retailer&   ")
	assert.Equal(t, 14, points)
}

func TestCalculatePointsIsRoundDollarAmount(t *testing.T) {
	points := calculatePointsIsRoundDollarAmount(100.00)
	assert.Equal(t, 50, points)
}

func TestCalculatePointsIsMultipleOf25(t *testing.T) {
	points := calculatePointsIsMultipleOf25(35.25)
	assert.Equal(t, 25, points)
}

func TestCalculatePointsForItems(t *testing.T) {
	items := []Item{{}, {}, {}, {}}
	points := calculatePointsForItems(items)
	assert.Equal(t, 10, points)
}

func TestCalculatePointsForItemDescription(t *testing.T) {
	items := []Item{
		{ShortDescription: "Item 1", Price: 10.00},
		{ShortDescription: "Item 2", Price: 12.00},
		{ShortDescription: "Apple", Price: 8.00},
	}
	points := calculatePointsForItemDescription(items)
	assert.Equal(t, 5, points)
}

func TestCalculatePointsForOddDay(t *testing.T) {
	points := calculatePointsForOddDay("2023-08-21")
	assert.Equal(t, 6, points)
}

func TestCalculatePointsForTimeOfPurchase(t *testing.T) {
	startTime, _ := time.Parse(HHMM, "14:00")
	endTime, _ := time.Parse(HHMM, "16:00")

	points := calculatePointsForTimeOfPurchase("15:00", startTime, endTime)
	assert.Equal(t, 10, points)
}
