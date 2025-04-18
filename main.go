package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/karan/watchlist/config"
	"github.com/karan/watchlist/routes"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Welcome to mongodb api")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config.ConnectDB()
	router := mux.NewRouter()
	routes.RegisterAllRoutes(router);
	log.Fatal(http.ListenAndServe(":4000", router))
}
