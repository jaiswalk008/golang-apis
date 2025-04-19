package routes

import (
	"github.com/gorilla/mux"
	"github.com/karan/watchlist/middleware"
)

// RegisterAllRoutes bundles all route registrations
func RegisterAllRoutes(r *mux.Router) {
	RegisterUserRoutes(r)
	authenticationRouter := r.PathPrefix("/").Subrouter()
	authenticationRouter.Use(middleware.AuthMiddleware)
	WatchlistRoutes(authenticationRouter)
	// Register more route groups as needed
}
