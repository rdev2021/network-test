package models

type PortCheckResponse struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Type         string `json:"type"`
	Status       string `json:"status"`
	ResponseTime int64  `json:"responseTime"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}
