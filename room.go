package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	ID           primitive.ObjectID `json:"-" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Owner        string             `json:"owner" bson:"owner"`
	UUID         string             `json:"uuid" bson:"uuid"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	Participants []string           `json:"participants" bson:"participants"`
}

func (r *Room) ToBSON() bson.D {
	return bson.D{
		{Key: "name", Value: r.Name},
		{Key: "owner", Value: r.Owner},
		{Key: "uuid", Value: r.UUID},
		{Key: "created_at", Value: r.CreatedAt},
		{Key: "participants", Value: r.Participants},
	}
}

type AddRoom struct {
	Name string `json:"name"`
}

type RemoveRoom struct {
	Name string `json:"name"`
}

type JoinRoom struct {
	Name string `json:"name"`
}

type LeaveRoom struct {
	Name string `json:"name"`
}

type GetRoom struct {
	Name string `json:"name"`
}
