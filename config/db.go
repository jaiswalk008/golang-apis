package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client // why as pointer and what is the use
/*
✅ 1. Why *mongo.Client (pointer) and not mongo.Client (value)?
Because mongo.Client is a struct that holds internal connection state, configuration, and manages connections to the MongoDB server. It’s quite a heavy object, and:

You don’t want to copy it around (copying structs is expensive and breaks state)

You want to share the same instance everywhere in your app

mongo.Connect() returns a *mongo.Client, not a plain mongo.Client
*/
func ConnectDB(){
	mongoURI := os.Getenv("MONGO_URI")
    ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client ,err := mongo.Connect(ctx , options.Client().ApplyURI(mongoURI))
	if err!=nil {
		log.Fatal("error in connecting to mongodb",err)
	}
	if err := client.Ping(ctx, nil); err != nil { // a ping test to check if the mongo db connection is successful or not
		log.Fatal("MongoDB ping failed:", err)
	}
	DB = client
	log.Println("MongoDb connected successfully")

}
func GetCollection(collectionName string) *mongo.Collection {
    return DB.Database("watchlist").Collection(collectionName)
}