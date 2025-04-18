package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
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

var (
	receipts = make(map[string]int)
	mutex    sync.Mutex
)

func processReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	points := calculatePoints(receipt)
	id := uuid.New().String()

	mutex.Lock()
	receipts[id] = points
	mutex.Unlock()

	response := map[string]string{"id": id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getPoints(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/receipts/")
	id = strings.TrimSuffix(id, "/points")

	mutex.Lock()
	points, exists := receipts[id]
	mutex.Unlock()

	if !exists {
		http.Error(w, "No receipt found for that ID.", http.StatusNotFound)
		return
	}

	response := map[string]int{"points": points}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func calculatePoints(r Receipt) int {
	points := 0

	// Rule 1: Alphanumeric characters in retailer name
	re := regexp.MustCompile(`[a-zA-Z0-9]`)
	points += len(re.FindAllString(r.Retailer, -1))

	// Rule 2 & 3: Total is round or multiple of 0.25
	total, _ := strconv.ParseFloat(r.Total, 64)
	if total == float64(int(total)) {
		points += 50
	}
	if math.Mod(total, 0.25) == 0 {
		points += 25
	}

	// Rule 4: 5 points per 2 items
	points += (len(r.Items) / 2) * 5

	// Rule 5: Length of item description is multiple of 3
	for _, item := range r.Items {
		trimmed := strings.TrimSpace(item.ShortDescription)
		if len(trimmed)%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(price * 0.2))
		}
	}

	// Rule 6: (Hidden rule: LLM-generated) 5 points if total > 10.00
	if total > 10.00 {
		points += 5
	}

	// Rule 7: Purchase day is odd
	if t, err := time.Parse("2006-01-02", r.PurchaseDate); err == nil && t.Day()%2 == 1 {
		points += 6
	}

	// Rule 8: Purchase time is after 2pm and before 4pm
	if t, err := time.Parse("15:04", r.PurchaseTime); err == nil {
		hour, min := t.Hour(), t.Minute()
		if (hour == 14 && min >= 0) || (hour == 15) {
			points += 10
		}
	}

	return points
}

func main() {
	http.HandleFunc("/receipts/process", processReceipt)
	http.HandleFunc("/receipts/", getPoints)
	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

