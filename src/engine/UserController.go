package engine

import (
	"context"
	"github.com/graph-gophers/graphql-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
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

func (d *DbEngine) Search(ctx context.Context, args struct{ Name string }) ([]*Person, error) {
	c := d.GetColl("msg")

	var objs []*Person
	re, err := c.Find(ctx, bson.M{"name": args.Name})
	if err != nil {
		return nil, err
	}
	err = re.All(ctx, &objs)
	if err != nil {
		return nil, err
	}

	return objs, nil
}

func (d *DbEngine) CreatePerson(ctx context.Context, args struct{ Input *InputPerson }) (*string, error) {
	c := d.GetColl("msg")

	p := Person{
		IdFiled:   primitive.NewObjectID(),
		Name:      args.Input.Name,
		Role:      args.Input.Role,
		Address:   args.Input.Address,
		Email:     args.Input.Email,
		Phone:     args.Input.Phone,
		CreatedAt: graphql.Time{time.Now().UTC()},
	}

	id := p.Id()

	_, err := c.InsertOne(context.Background(), &p)
	if err != nil {
		return &id, err
	}

	return &id, nil
}
