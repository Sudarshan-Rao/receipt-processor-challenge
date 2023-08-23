package main

import (
	"math"
	"strings"
	"time"
	"unicode"
)

func CalculateTotalPoints(receipt Receipt) int {
	points := 0

	points += calculatePointsIsAlphanumeric(receipt.Retailer)
	points += calculatePointsIsRoundDollarAmount(receipt.Total)
	points += calculatePointsIsMultipleOf25(receipt.Total)
	points += calculatePointsForItems(receipt.Items)
	points += calculatePointsForItemDescription(receipt.Items)
	points += calculatePointsForOddDay(receipt.PurchaseDate)
	points += calculatePointsForTimeOfPurchase(receipt.PurchaseTime, STARTTIME2PM, ENDTIME4PM)

	return points
}

func calculatePointsIsAlphanumeric(retailer string) int {
	points := 0
	pointValue := 1

	for _, char := range retailer {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			points += pointValue
		}
	}

	return points
}

func calculatePointsIsRoundDollarAmount(total float32) int {
	points := 0
	pointValue := 50

	totalDollars := int(total)
	if float32(totalDollars) == total {
		points += pointValue
	}

	return points
}

func calculatePointsIsMultipleOf25(total float32) int {
	points := 0
	pointValue := 25

	if math.Mod(float64(total), 0.25) == 0 {
		points += pointValue
	}

	return points
}

func calculatePointsForItems(items []Item) int {
	points := 0
	pointValue := 5

	points += (len(items) / 2) * pointValue

	return points
}

func calculatePointsForItemDescription(items []Item) int {
	points := 0
	multipleOf := 3
	multiplyPriceBy := 0.2

	for _, item := range items {
		trimmedLen := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLen%multipleOf == 0 {
			itemPoints := int(math.Ceil(float64(item.Price) * multiplyPriceBy))
			points += itemPoints
		}
	}

	return points
}

func calculatePointsForOddDay(purchaseDate string) int {
	points := 0
	pointValue := 6

	parsedPurchaseDate, err := time.Parse(YYYY_MM_DD, purchaseDate)
	if err == nil && parsedPurchaseDate.Day()%2 != 0 {
		points += pointValue
	}

	return points
}

func calculatePointsForTimeOfPurchase(purchaseTime string, startTime time.Time, endTime time.Time) int {
	points := 0
	pointValue := 10

	parsedPurchaseTime, err := time.Parse(HHMM, purchaseTime)
	if err == nil {
		if parsedPurchaseTime.After(startTime) && parsedPurchaseTime.Before(endTime) {
			points += pointValue
		}
	}

	return points
}
