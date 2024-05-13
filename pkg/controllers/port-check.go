package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rdev2021/network-test/pkg/models"
	"github.com/rdev2021/network-test/pkg/utils"
)

const ContentType = "Content-Type"
const JsonContentType = "application/json"

type PortCheckRequst struct {
	HostName string `json:"hostname"`
	Port     string `json:"port"`
}

func CheckPortHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path, r.RemoteAddr)
	start := time.Now()

	var data PortCheckRequst

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		responseTime := time.Since(start).Milliseconds()
		result := models.PortCheckResponse{
			Status:       "Error",
			ResponseTime: responseTime,
			ErrorMessage: fmt.Sprintf("Error decoding JSON: %v", err),
		}
		w.Header().Set(ContentType, JsonContentType)
		json.NewEncoder(w).Encode(result)
		return
	}

	port, err := strconv.Atoi(strings.TrimSpace(data.Port))
	if err != nil {
		responseTime := time.Since(start).Milliseconds()
		result := models.PortCheckResponse{
			Host:         data.HostName,
			Port:         0,
			Status:       "Error",
			ResponseTime: responseTime,
			ErrorMessage: "Invalid port number: " + err.Error(),
		}

		w.Header().Set(ContentType, JsonContentType)
		json.NewEncoder(w).Encode(result)
		return
	}

	result := utils.PortCheck(strings.TrimSpace(data.HostName), port, "TCP")

	w.Header().Set(ContentType, JsonContentType)
	json.NewEncoder(w).Encode(result)
}
