package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/karan/watchlist/config"
	"github.com/karan/watchlist/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	var user model.User
	//parse request body
	err := json.NewDecoder(r.Body).Decode(&user)
	/*
	r.Body - This is the raw request body from the HTTP request, which contains the JSON data sent by the client [1]

json.NewDecoder(r.Body) - Creates a new decoder that reads from the request body stream

.Decode(&user) - Reads the JSON data and converts it into your Go struct (model.User in this case)
*/
	if err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
	}
	//checking for existing collection
	userCollection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser model.User
	
	if err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser); err == nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	} else if err != mongo.ErrNoDocuments {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error in hashing password", http.StatusInternalServerError)
	}
	user.ID = primitive.NewObjectID()
	user.Password = string(hashedPassword)
	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		fmt.Println("Insert error:", err)
		http.Error(w, "Error saving user", http.StatusInternalServerError)
		return
	}

	fmt.Println("user: ")
	fmt.Println(user)

	

	//Response
	response := map[string]any{
		"message": "user created successfully",
		"user": map[string]any{
			"id":    user.ID.Hex(),
			"name":  user.Name,
			"email": user.Email,
		},
	}
	fmt.Println("response: ")
	fmt.Println(response)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
