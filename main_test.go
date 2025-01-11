package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProcessReceipt1(t *testing.T) {
	//Sets up the router with routes
	r := Routes()

	//The test receipt in JSON form
	receipt := `{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}`
	//Opens a new HTTP request
	req := httptest.NewRequest(http.MethodPost, "/receipts/process", bytes.NewBuffer([]byte(receipt)))
	req.Header.Set("Content-Type", "application/json")

	//Responce Recorder to record results
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	//Checks status condition of the last request
	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("expected status %d but got %d", http.StatusOK, status)
	}

	//Parsing through the JSON responce
	var response map[string]string
	err2 := json.Unmarshal(rr.Body.Bytes(), &response)
	if err2 != nil {
		t.Fatalf("failed to parse response: %v", err2)
	}

	//Checks to see if an ID was returned in the JSON responce
	id, ok := response["id"]
	if !ok {
		t.Fatalf("response did not contain an 'id' field")
	}

	//Testing the get points endpoint
	pointsReq := httptest.NewRequest(http.MethodGet, "/receipts/"+id+"/points", nil)
	pointsRR := httptest.NewRecorder()
	r.ServeHTTP(pointsRR, pointsReq)

	//Checks status condition of the last request
	status2 := pointsRR.Code
	if status2 != http.StatusOK {
		t.Errorf("expected status %d but got %d", http.StatusOK, status2)
	}

	//Parses the points responce
	var pointsResponse map[string]int
	err3 := json.Unmarshal(pointsRR.Body.Bytes(), &pointsResponse)
	if err3 != nil {
		t.Fatalf("failed to parse points response: %v", err3)
	}

	//Checks if the expected amount of points matches with what was returned in the JSON responce
	expectedPoints := 28
	points, ok := pointsResponse["points"]
	if !ok || points != expectedPoints {
		t.Errorf("expected points %d but got %d", expectedPoints, points)
	}
}

func TestProcessReceipt2(t *testing.T) {
	//Sets up the router with routes
	r := Routes()

	//The test receipt in JSON form
	receipt := `{
  "retailer": "M&M Corner Market",
  "purchaseDate": "2022-03-20",
  "purchaseTime": "14:33",
  "items": [
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    }
  ],
  "total": "9.00"
}`
	//Opens a new HTTP request
	req := httptest.NewRequest(http.MethodPost, "/receipts/process", bytes.NewBuffer([]byte(receipt)))
	req.Header.Set("Content-Type", "application/json")

	//Responce Recorder to record results
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	//Checks status condition of the last request
	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("expected status %d but got %d", http.StatusOK, status)
	}

	//Parsing through the JSON responce
	var response map[string]string
	err2 := json.Unmarshal(rr.Body.Bytes(), &response)
	if err2 != nil {
		t.Fatalf("failed to parse response: %v", err2)
	}

	//Checks to see if an ID was returned in the JSON responce
	id, ok := response["id"]
	if !ok {
		t.Fatalf("response did not contain an 'id' field")
	}

	//Testing the get points endpoint
	pointsReq := httptest.NewRequest(http.MethodGet, "/receipts/"+id+"/points", nil)
	pointsRR := httptest.NewRecorder()
	r.ServeHTTP(pointsRR, pointsReq)

	//Checks status condition of the last request
	status2 := pointsRR.Code
	if status2 != http.StatusOK {
		t.Errorf("expected status %d but got %d", http.StatusOK, status2)
	}

	//Parses the points responce
	var pointsResponse map[string]int
	err3 := json.Unmarshal(pointsRR.Body.Bytes(), &pointsResponse)
	if err3 != nil {
		t.Fatalf("failed to parse points response: %v", err3)
	}

	//Checks if the expected amount of points matches with what was returned in the JSON responce
	expectedPoints := 109
	points, ok := pointsResponse["points"]
	if !ok || points != expectedPoints {
		t.Errorf("expected points %d but got %d", expectedPoints, points)
	}
}
