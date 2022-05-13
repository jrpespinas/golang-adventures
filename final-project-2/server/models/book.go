package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	User_ID    primitive.ObjectID `json:"user_id,omitempty" bson:"user_id, omitempty"`
	Title      string             `json:"title,omitempty" bson:"title, omitempty"`
	Author     string             `json:"author,omitempty" bson:"author, omitempty"`
	IsFinished bool               `json:"is_finished,omitempty" bson:"is_finished, omitempty"`
}
