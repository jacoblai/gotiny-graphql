package models

import (
	"github.com/graph-gophers/graphql-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Person struct {
	IdFiled   primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Role      string             `json:"role" bson:"role"`
	Address   *[]string          `json:"address" bson:"address"`
	CreatedAt graphql.Time       `json:"created_at" bson:"created_at"`
	Email     string             `json:"email" bson:"email"`
	Phone     string             `json:"phone" bson:"phone"`
}

func (u Person) Id() string {
	return u.IdFiled.Hex()
}

type InputPerson struct {
	Name      string
	Role      string
	Address   *[]string
	CreatedAt *graphql.Time
	Email     string
	Phone     string
}
