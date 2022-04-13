package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Username string             `json:"username" bson:"username"`
}

func (u User) ToBSON() bson.D {
	return bson.D{
		{Key: "username", Value: u.Username},
	}
}

type AddUser struct {
	Username string `json:"username"`
}

type GetUser struct {
	Username string `json:"username"`
}
