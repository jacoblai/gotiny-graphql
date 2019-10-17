package models

import (
	"github.com/graph-gophers/graphql-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var T_Person = "persons"

type Person struct {
	IdFiled        primitive.ObjectID `json:"_id" bson:"_id,omitempty" jsonschema:"-"`
	Name           string             `json:"name,omitempty" bson:"name,omitempty" jsonschema:"required,minLength=2,maxLength=64"`
	Role           string             `json:"role,omitempty" bson:"role,omitempty" jsonschema:"enum=ADMIN|USER"`
	Address        *[]string          `json:"address,omitempty" bson:"address,omitempty" jsonschema:"-"`
	CreatedAtFiled time.Time          `json:"createat,omitempty" bson:"createat,omitempty"`
	Email          string             `json:"email,omitempty" bson:"email,omitempty" jsonschema:"required,minLength=2,maxLength=64"`
	Phone          string             `json:"phone,omitempty" bson:"phone,omitempty" jsonschema:"required,minLength=2,maxLength=64"`
	Total          *float64           `json:"total,omitempty" bson:"total,omitempty"`
	Order          *int32             `json:"order,omitempty" bson:"order,omitempty"`
}

func (u Person) CreatedAt() graphql.Time {
	return graphql.Time{Time: u.CreatedAtFiled.Local()}
}

func (u Person) Id() string {
	return u.IdFiled.Hex()
}

type InputPerson struct {
	Name    string
	Email   string
	Role    string
	Phone   string
	Address *[]string
	Total   *float64
	Order   *int32
}
