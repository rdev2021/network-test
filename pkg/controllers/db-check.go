package controllers

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/rdev2021/network-test/pkg/utils"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"

	// _ "github.com/godror/godror"
	_ "github.com/lib/pq"
)

type DBCheckRequest struct {
	DbType       string `json:"dbType"`
	Hostname     string `json:"hostname"`
	Port         string `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	ServiceName  string `json:"serviceName,omitempty"`
	DatabaseName string `json:"databaseName,omitempty"`
}

type DBCheckResponse struct {
	Host         string `json:"host,omitempty"`
	Port         int    `json:"port,omitempty"`
	UserName     string `json:"userName,omitempty"`
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	ResponseTime int64  `json:"responseTime"`
}

func CheckDBHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path, r.RemoteAddr)
	start := time.Now()

	userData := r.Header.Get("userdata")

	decoded, err := base64.StdEncoding.DecodeString(userData)
	if err != nil {
		responseTime := time.Since(start).Milliseconds()
		result := DBCheckResponse{
			Status:       "Error",
			ResponseTime: responseTime,
			ErrorMessage: err.Error(),
		}
		w.Header().Set(ContentType, JsonContentType)
		json.NewEncoder(w).Encode(result)

		return
	}

	var data DBCheckRequest
	err = json.Unmarshal(decoded, &data)
	if err != nil {
		responseTime := time.Since(start).Milliseconds()
		result := DBCheckResponse{
			Status:       "Error",
			ResponseTime: responseTime,
			ErrorMessage: err.Error(),
		}
		w.Header().Set(ContentType, JsonContentType)
		json.NewEncoder(w).Encode(result)

		return
	}

	portInt, err := strconv.Atoi(data.Port)
	if err != nil {
		responseTime := time.Since(start).Milliseconds()
		result := DBCheckResponse{
			Status:       "Error",
			ResponseTime: responseTime,
			ErrorMessage: "Invalid port number: " + err.Error(),
		}
		w.Header().Set(ContentType, JsonContentType)
		json.NewEncoder(w).Encode(result)

		return
	}
	result := utils.PortCheck(data.Hostname, portInt, "TCP")

	if result.Status != "Connected" {
		responseTime := time.Since(start).Milliseconds()
		result := DBCheckResponse{
			Host:         data.Hostname,
			Port:         portInt,
			Status:       "Error",
			ResponseTime: responseTime,
			ErrorMessage: result.ErrorMessage,
		}
		w.Header().Set(ContentType, JsonContentType)
		json.NewEncoder(w).Encode(result)
		return
	}

	dsn := dsn(data)

	db, err := sql.Open(data.DbType, dsn)
	if err != nil {
		responseTime := time.Since(start).Milliseconds()
		result := DBCheckResponse{
			Status:       "Failure",
			ErrorMessage: err.Error(),
			ResponseTime: responseTime,
		}
		w.Header().Set(ContentType, JsonContentType)
		json.NewEncoder(w).Encode(result)
		return
	}

	err = db.Ping() // Test connection

	if err != nil {

		responseTime := time.Since(start).Milliseconds()
		result := DBCheckResponse{
			Status:       "Failure",
			ErrorMessage: err.Error(),
			ResponseTime: responseTime,
		}
		w.Header().Set(ContentType, JsonContentType)
		json.NewEncoder(w).Encode(result)
		return
	}
	defer db.Close()

	responseTime := time.Since(start).Milliseconds()
	response := DBCheckResponse{
		Host:         data.Hostname,
		Port:         portInt,
		UserName:     data.Username,
		Status:       "Success",
		ResponseTime: responseTime,
	}

	w.Header().Set(ContentType, JsonContentType)
	json.NewEncoder(w).Encode(response)
}

func dsn(data DBCheckRequest) string {
	switch data.DbType {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/", data.Username, data.Password, data.Hostname, data.Port)
	case "postgres":
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", data.Hostname, data.Port, data.Username, data.Password, data.DatabaseName)
	case "mssql":
		return fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s", data.Hostname, data.Username, data.Password, data.Port)
	// case "oracle":
	// 	return fmt.Sprintf("%s/%s@//%s:%s/%s", data.Username, data.Password, data.Hostname, data.Port, data.ServiceName)
	default:
		return ""
	}
}
