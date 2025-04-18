package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Watchlist struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
	MovieName string `bson:"movie_name,omitempty" json:"movie_name,omitempty"`
	Watched bool `bson:"watched" json:"watched"`
}