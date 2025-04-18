package routes

import (
	"github.com/gorilla/mux"
	"github.com/karan/watchlist/controller"
)

func RegisterUserRoutes(r *mux.Router){
	r.HandleFunc("/signup",controller.SignupHandler).Methods("POST")
}