package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Recipe struct {
	ID           primitive.ObjectID `json:"id"           bson:"_id"`
	Name         string             `json:"name"         bson:"name"`
	Tags         []string           `json:"tags"         bson:"tags"`
	Ingredients  []string           `json:"ingredients"  bson:"ingredients"`
	Instructions []string           `json:"instructions" bson:"instructions"`
	PublishedAt  time.Time          `json:"publishedAt"  bson:"publishedAt"`
}

type User struct {
	Password string `json:"password"`
	Username string `json:"username"`
}
