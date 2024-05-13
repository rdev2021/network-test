package routes

import (
	"github.com/gorilla/mux"
	"github.com/rdev2021/network-test/pkg/controllers"
)

var RegisterCCCRoutes = func(router *mux.Router) {
	router.HandleFunc("/port", controllers.CheckPortHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/resolve", controllers.CheckDomainHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/http", controllers.CheckHttpHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/db", controllers.CheckDBHandler).Methods("POST", "OPTIONS")
}
