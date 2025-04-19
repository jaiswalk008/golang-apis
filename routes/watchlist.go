package routes

import (
	"github.com/gorilla/mux"
	"github.com/karan/watchlist/controller"
)

func WatchlistRoutes(r *mux.Router) {
	r.HandleFunc("/watchlist",controller.AddWatchlistHandler).Methods("POST")
	r.HandleFunc("/watchlist",controller.GetAllWatchlist).Methods("GET")
	r.HandleFunc("/watchlist/{id}",controller.UpdateWatchList).Methods("PATCH")
	r.HandleFunc("/watchlist/{id}",controller.DeleteWatchlist).Methods("DELETE")
}