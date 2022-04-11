package main

import "time"

type Question struct {
	Text      string    `json:"question" bson:"question"`
	Active    bool      `json:"active" bson:"active"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
