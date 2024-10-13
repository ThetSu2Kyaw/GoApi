package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Define a struct that matches the structure of the API response
type FlightAPIResponse struct {
	AeroCRS struct {
		Success bool `json:"success"`
		Flights struct {
			Count  int           `json:"count"`
			Flight []interface{} `json:"flight"`
		} `json:"flights"`
	} `json:"aerocrs"`
}

func fetchFlightData() (*FlightAPIResponse, error) {
	url := "https://api.aerocrs.com/v5/getDeepLink?from=RGN&to=MDL&start=2024/08/30&adults=1&child=0&infant=0&chargetype=NTL&currency=MMK"

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"auth_id":       {"123456"},
		"auth_password": {"password"},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var flightData FlightAPIResponse
	err = json.Unmarshal(data, &flightData)
	if err != nil {
		return nil, err
	}

	return &flightData, nil
}

func flightsHandler(w http.ResponseWriter, r *http.Request) {
	data, err := fetchFlightData()
	if err != nil {
		http.Error(w, "Failed to fetch flight data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	json.NewEncoder(w).Encode(data)
}

func main() {
	http.HandleFunc("/flights", flightsHandler)
	fmt.Println("Server is running on http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}
