package engine

import (
	"context"
	"github.com/graph-gophers/graphql-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"models"
	"time"
)

func (d *DbEngine) Search(ctx context.Context, args struct{ Name string }) ([]*models.Person, error) {
	c := d.GetColl("msg")

	var objs []*models.Person
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

func (d *DbEngine) CreatePerson(ctx context.Context, args struct{ Input *models.InputPerson }) (*string, error) {
	c := d.GetColl("msg")

	p := models.Person{
		IdFiled:   primitive.NewObjectID(),
		Name:      args.Input.Name,
		Role:      args.Input.Role,
		Address:   args.Input.Address,
		Email:     args.Input.Email,
		Phone:     args.Input.Phone,
		CreatedAt: graphql.Time{time.Now().Local()},
	}

	id := p.Id()

	_, err := c.InsertOne(context.Background(), &p)
	if err != nil {
		return &id, err
	}

	return &id, nil
}
