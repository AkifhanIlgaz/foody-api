package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Session struct {
	Uid       primitive.ObjectID `bson:"uid"`
	Token     string             `bson:"-"`
	TokenHash string             `bson:"tokenHash"`
}
