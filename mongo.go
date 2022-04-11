package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName = "chat"

type mongoDB struct {
	client          *mongo.Client
	usersCollection *mongo.Collection
	roomsCollection *mongo.Collection
}

func NewMongoDB(ctx context.Context, host string, port int) (*mongoDB, error) {
	// TODO remove creds
	client, err := mongo.Connect(
		ctx,
		options.Client().
			ApplyURI(fmt.Sprintf("mongodb://%s:%d", host, port)).
			SetAuth(options.Credential{Username: "admin", Password: "secret"}))
	if err != nil {
		return nil, err
	}

	// collections
	usersCollection := client.Database(dbName).Collection("users")
	roomsCollection := client.Database(dbName).Collection("rooms")

	// idxes
	_, err = roomsCollection.Indexes().CreateMany(
		ctx,
		[]mongo.IndexModel{
			{
				Keys:    bson.M{"name": 1},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys:    bson.M{"owner": 1},
				Options: nil,
			},
		})
	if err != nil {
		return nil, err
	}

	_, err = usersCollection.Indexes().CreateMany(
		ctx,
		[]mongo.IndexModel{
			{
				Keys:    bson.M{"username": 1},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys:    bson.M{"token": 1},
				Options: options.Index().SetUnique(true),
			},
		})
	if err != nil {
		return nil, err
	}

	return &mongoDB{client, usersCollection, roomsCollection}, nil
}

func (m *mongoDB) Ping(ctx context.Context) error {
	return m.client.Ping(ctx, nil)
}

func (m *mongoDB) AddUser(ctx context.Context, user *User) error {
	_, err := m.usersCollection.InsertOne(ctx, user.ToBSON())
	return err
}

func (m *mongoDB) GetUser(ctx context.Context, username, token string) (*User, error) {
	filters := make(bson.D, 0, 2)
	if username != "" {
		filters = append(filters, bson.E{Key: "username", Value: username})
	}
	if token != "" {
		filters = append(filters, bson.E{Key: "token", Value: token})
	}

	res := m.usersCollection.FindOne(ctx, filters)
	if err := res.Err(); err != nil {
		return nil, err
	}

	var user User
	if err := res.Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *mongoDB) UserExists(ctx context.Context, username string) (bool, error) {
	_, err := m.GetUser(ctx, username, "")
	if err == nil {
		return true, nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}

	return false, err
}

func (m *mongoDB) AddRoom(ctx context.Context, name string, user *User) (*Room, error) {
	room := Room{
		UUID:         uuid.NewString(),
		Name:         name,
		Owner:        user.Username,
		CreatedAt:    time.Now().In(time.UTC),
		Participants: []string{},
	}
	_, err := m.roomsCollection.InsertOne(ctx, room.ToBSON())
	if err != nil {
		return nil, err
	}

	return &room, nil
}

func (m *mongoDB) RemoveRoom(ctx context.Context, name string, user *User) error {
	_, err := m.roomsCollection.DeleteOne(
		ctx,
		bson.D{
			{Key: "owner", Value: user.Username},
			{Key: "name", Value: name},
		})
	return err
}

func (m *mongoDB) JoinRoom(ctx context.Context, roomName string, user *User) (*Room, error) {
	// add user to participants if they are not already there
	res := m.roomsCollection.FindOneAndUpdate(
		ctx,
		bson.M{"name": roomName, "participants": bson.M{"$nin": []string{user.Username}}},
		bson.M{"$push": bson.M{"participants": user.Username}})
	if err := res.Err(); err != nil {
		return nil, err
	}

	var room *Room
	if err := res.Decode(&room); err != nil {
		return nil, err
	}

	room.Participants = append(room.Participants, user.Username)
	return room, nil
}

func (m *mongoDB) LeaveRoom(ctx context.Context, roomName string, user *User) error {
	res := m.roomsCollection.FindOneAndUpdate(
		ctx,
		bson.D{{Key: "name", Value: roomName}},
		bson.M{"$pull": bson.M{"participants": user.Username}})

	return res.Err()
}

func (m *mongoDB) GetRoom(ctx context.Context, name string) (*Room, error) {
	res := m.roomsCollection.FindOne(ctx, bson.M{"name": name})
	if err := res.Err(); err != nil {
		return nil, err
	}

	var room *Room
	if err := res.Decode(&room); err != nil {
		return nil, err
	}

	return room, nil
}

func (m *mongoDB) GetUserRooms(ctx context.Context, user *User) ([]*Room, error) {
	cur, err := m.roomsCollection.Find(ctx, bson.M{"owner": user.Username})
	if err != nil {
		return nil, err
	}

	var rooms []*Room
	for cur.Next(ctx) {
		var room *Room
		if err := cur.Decode(&room); err != nil {
			continue
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}
