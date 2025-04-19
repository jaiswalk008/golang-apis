package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/karan/watchlist/config"
	"github.com/karan/watchlist/middleware"
	"github.com/karan/watchlist/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddWatchlistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only post method allowed", http.StatusMethodNotAllowed)
		return
	}
	var watchlist, existingWatchlist model.Watchlist
	if err := json.NewDecoder(r.Body).Decode(&watchlist); err != nil {
		http.Error(w, "json parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	watchlistCollection := config.GetCollection("watchlist")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := watchlistCollection.FindOne(ctx, bson.M{"movie_name": watchlist.MovieName}).Decode(&existingWatchlist); err == nil {
		fmt.Println(err)
		http.Error(w, "Movie already in the wartchlist", http.StatusConflict)
		return;
	}
	watchlist.ID = primitive.NewObjectID()
	userId := r.Context().Value(middleware.UserIDKey).(string)
	userObjectId,err := primitive.ObjectIDFromHex(userId)
	if err!= nil{
		http.Error(w, "Invalid user ID format"+err.Error(), http.StatusBadRequest)
		return
	}
	watchlist.UserID = userObjectId
	response := map[string]any {
		"message":"Watchlist added successfully",
		"watchlist" : map[string]any {
			"id":watchlist.ID.Hex(),
			"movieName":watchlist.MovieName,
			"watched":watchlist.Watched,
		},
	}
	json.NewEncoder(w).Encode(response);


}
