package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}
/*
bson:"_id,omitempty":

This means that the ID field will be mapped to MongoDB's _id field.

The omitempty directive means that if ID is empty (i.e., nil or zero value), it will be ignored when marshaling the object (so it won’t be written to the database).
*/