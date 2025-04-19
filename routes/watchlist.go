package routes

import (
	"github.com/gorilla/mux"
	"github.com/karan/watchlist/controller"
)

func WatchlistRoutes(r *mux.Router) {
	r.HandleFunc("/watchlist",controller.AddWatchlistHandler).Methods("POST")
}