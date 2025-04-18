package routes

import "github.com/gorilla/mux"

// RegisterAllRoutes bundles all route registrations
func RegisterAllRoutes(r *mux.Router) {
	RegisterUserRoutes(r)
	// RegisterWatchlistRoutes(r)
	// Register more route groups as needed
}
