package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// Maps the ID of each reciept with its points
var receiptData = make(map[string]int)

// This function sets up the router and endpoints to redirect to its assigned method
func Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Route("/receipts/process", func(r chi.Router) {
		r.Post("/", processReceiptHandler)
	})

	router.Route("/receipts/{id}/points", func(r chi.Router) {
		r.Get("/", getPointsHandler)
	})

	return router
}

// handles the POST request to process a receipt
func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt

	//Parses the JSON file
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid receipt JSON", http.StatusBadRequest)
		return
	}

	//Calls method to create an ID
	id := generateID(receipt)

	//Calls a method to calculate the amount of points earned
	points := calculatePoints(receipt)

	//Stores the ID and points pair
	receiptData[id] = points

	//Returns the ID as JSON
	response := map[string]string{"id": id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handles the GET request to fetch points for a receipt
func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	//Checks to see if the ID exists
	points, exists := receiptData[id]
	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	//Returns the points in a JSON format
	response := map[string]int{"points": points}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// generateID generates a unique ID for a receipt using the SHA256 hashing function
func generateID(receipt Receipt) string {
	data := receipt.Retailer + receipt.PurchaseDate + receipt.PurchaseTime + receipt.Total
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// calculatePoints calculates points based on the rules
func calculatePoints(receipt Receipt) int {
	points := 0

	//1pt for each alphanumeric character in the store name
	reg := regexp.MustCompile(`[a-zA-Z0-9]`)
	points += len(reg.FindAllString(receipt.Retailer, -1))

	//50pts if the total ends in a round dollar amount
	total, err := strconv.ParseFloat(receipt.Total, 64) //begins by converting the total into a float
	if err != nil {                                     //check to see if there is no error as a result
		log.Print("Error in parsing to float", err)
	}
	if math.Mod(total, 1.0) == 0 { //does final calculation to ensure it is a whole number
		points += 50
	}

	//25pts if the total is a mutiple of ".25"
	total25, err2 := strconv.ParseFloat(receipt.Total, 64) //begins by converting the total into a float
	if err2 != nil {                                       //check to see if there is no error as a result
		log.Print("Error in parsing to float", err2)
	}
	if math.Mod(total25, .25) == 0 { //does final calculation to ensure it is a whole number
		points += 25
	}

	//5pts for every two items on the reciept
	points += int(len(receipt.Items)/2) * 5

	//If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	for _, item := range receipt.Items {
		descriptionLength := len(strings.TrimSpace(item.ShortDescription)) //removes all whitespace
		if descriptionLength%3 == 0 && descriptionLength > 0 {
			price, err3 := strconv.ParseFloat(item.Price, 64)
			if err3 != nil { //check to see if there is no error as a result
				log.Print("Error in parsing to float price of item", err2)
			}
			points += int(math.Ceil(price * 0.2)) //calculates 20% of the item price, then adds rounds to the nearest integer to add (ex: 1.6 -> 2)

		}
	}

	//6pts if the day of purchase is odd
	date, err4 := time.Parse("2006-01-02", receipt.PurchaseDate) //uses YYYY-MM-DD format
	if err4 != nil {
		log.Print("Error in parsing date of purchase", err4)
	}
	if date.Day()%2 != 0 { //modulus to ensure its an odd date
		points += 6
	}

	//10pts if the time of purchase is between 2pm and 4pm
	purchaseTime, err5 := time.Parse("15:04", receipt.PurchaseTime) //24 hour clock format
	if err5 != nil {
		log.Print("Error in parsing time of purchase", err5)
	}
	if purchaseTime.Hour() >= 14 && purchaseTime.Hour() <= 16 {
		points += 10
	}

	return points
}
