package controllers

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type HttpCheckRequest struct {
	URL string `json:"url"`
}

type HttpCheckResponse struct {
	URL          string `json:"url"`
	Success      bool   `json:"success"`
	StatusCode   int    `json:"statusCode,omitempty"`
	ResponseTime int64  `json:"responseTime"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

func CheckHttpHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path, r.RemoteAddr)
	start := time.Now()

	var data HttpCheckRequest
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		responseTime := time.Since(start).Milliseconds()
		result := HttpCheckResponse{
			Success:      false,
			ResponseTime: responseTime,
			ErrorMessage: fmt.Sprintf("Error decoding JSON: %v", err),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		return
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Make HTTP GET request
	resp, err := client.Get(data.URL)
	if err != nil {
		responseTime := time.Since(start).Milliseconds()
		result := HttpCheckResponse{
			URL:          data.URL,
			Success:      false,
			ResponseTime: responseTime,
			ErrorMessage: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		return
	}

	defer resp.Body.Close()

	responseTime := time.Since(start).Milliseconds()
	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	result := HttpCheckResponse{
		URL:          data.URL,
		Success:      success,
		StatusCode:   resp.StatusCode,
		ResponseTime: responseTime,
	}

	// Encode response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
