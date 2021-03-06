package main

import (
	"context"
)

type DB interface {
	Ping(ctx context.Context) error
	AddUser(ctx context.Context, user *User) error
	UserExists(ctx context.Context, username string) (bool, error)
	AddRoom(ctx context.Context, name string, user *User) (*Room, error)
	RemoveRoom(ctx context.Context, name string, user *User) error
	JoinRoom(ctx context.Context, roomName string, user *User) (*Room, error)
	LeaveRoom(ctx context.Context, roomName string, user *User) error
	GetRoom(ctx context.Context, name string) (*Room, error)
	GetUser(ctx context.Context, username string) (*User, error)
	GetUserRooms(ctx context.Context, user *User) ([]*Room, error)
}
