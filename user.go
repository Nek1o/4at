package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Username string             `json:"username" bson:"username"`
	Token    string             `json:"token" bson:"token"`
}

func (u User) ToBSON() bson.D {
	return bson.D{
		{Key: "username", Value: u.Username},
		{Key: "token", Value: u.Token},
	}
}

type AddUser struct {
	Username string `json:"username"`
}

type GetUser struct {
	Username string `json:"username"`
}
