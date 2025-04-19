package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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
	var watchlistCollection = config.GetCollection("watchlist")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := watchlistCollection.FindOne(ctx, bson.M{"movie_name": watchlist.MovieName}).Decode(&existingWatchlist); err == nil {
		fmt.Println(err)
		http.Error(w, "Movie already in the wartchlist", http.StatusConflict)
		return
	}
	watchlist.ID = primitive.NewObjectID()
	userId := r.Context().Value(middleware.UserIDKey).(string)
	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		http.Error(w, "Invalid user ID format"+err.Error(), http.StatusBadRequest)
		return
	}
	watchlist.UserID = userObjectId

	// Insert the watchlist into the database
	_, err = watchlistCollection.InsertOne(ctx, watchlist)
	if err != nil {
		http.Error(w, "Error saving watchlist: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"message": "Watchlist added successfully",
		"watchlist": map[string]any{
			"id":        watchlist.ID.Hex(),
			"movieName": watchlist.MovieName,
			"watched":   watchlist.Watched,
		},
	}
	json.NewEncoder(w).Encode(response)

}

func GetAllWatchlist(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}
	userId := r.Context().Value(middleware.UserIDKey).(string)
	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusBadRequest)
	}
	var watchlistCollection = config.GetCollection("watchlist")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := watchlistCollection.Find(ctx, bson.M{"user_id": userObjectId})
	if err != nil {
		http.Error(w, "Error in fetching the watchlist : "+err.Error(), http.StatusInternalServerError)
	}
	defer cursor.Close(ctx)

	/*
		In MongoDB (and Go's mongo-go-driver):

		When you use .Find(), MongoDB doesn't return all data at once.

		It gives you a cursor — like a pointer that you iterate over to read documents from the result set.

		You can use:

		.Next() to loop manually,

		or .All() to decode all documents into a slice (like you're doing).
	*/

	var watchlists []model.Watchlist
	if err = cursor.All(ctx, &watchlists); err != nil {
		http.Error(w, "Error decoding watchlists: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	// Send the JSON response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Watchlists fetched successfully",
		"watchlists": watchlists,
	})
}
func UpdateWatchList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		http.Error(w, "Only PATCH method is required", http.StatusMethodNotAllowed)
		return
	}
	params := mux.Vars(r)
	watchlistId := params["id"]
	fmt.Println("id, = ",watchlistId)
	watchlistObjectId, _ := primitive.ObjectIDFromHex(watchlistId)
	var updatedWatchList model.Watchlist
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := json.NewDecoder(r.Body).Decode(&updatedWatchList); err != nil {
		http.Error(w, "error in json data", http.StatusInternalServerError)
		return
	}
	watchlistCollection := config.GetCollection("watchlist")

	result, err := watchlistCollection.UpdateByID(ctx, watchlistObjectId, bson.M{"$set": updatedWatchList})
	fmt.Println(result)
	if err != nil {
		http.Error(w, "Error in updating the watchlist: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if result.MatchedCount == 0 {
		http.Error(w, "No watchlist found with that ID", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Watchlist updated successfully",
	})

}
func DeleteWatchlist (w http.ResponseWriter , r *http.Request){
	params := mux.Vars(r)
	watchlistId := params["id"]
	watchlistObjectId,_:= primitive.ObjectIDFromHex(watchlistId)
	ctx,cancel := context.WithTimeout(context.Background() , 10*time.Second)
	defer cancel()
	watchlistCollection := config.GetCollection("watchlist")
	result ,err := watchlistCollection.DeleteOne(ctx,bson.M{"_id":watchlistObjectId})
	if err!=nil {
		http.Error(w,"Failed to delete the watchlist",http.StatusInternalServerError)
		return 
	}
	if result.DeletedCount==0 {
		http.Error(w,"No watchlist to delete",http.StatusNotFound)
		return
	}
	response:= map[string]string {
		"message":"Watchlist deleted successfully",
	}
	json.NewEncoder(w).Encode(response)
}