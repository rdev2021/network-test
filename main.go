package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rdev2021/network-test/pkg/routes"
)

func main() {
	myRouter := mux.NewRouter()
	routes.RegisterCCCRoutes(myRouter)

	bs := http.FileServer(http.Dir("./ui/"))
	myRouter.PathPrefix("/home/").Handler(http.StripPrefix("/home/", bs))

	srv := &http.Server{
		Handler:      myRouter,
		Addr:         ":9091",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Server running on port 9091...")
	log.Fatal(srv.ListenAndServe())
}
