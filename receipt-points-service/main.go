package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func initializeRouter() {
	r := mux.NewRouter()

	r.HandleFunc("/receipts/process", processReceipt).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", getReceiptPoints).Methods("GET")

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":9000", r))
}
func main() {
	initializeRouter()
}
