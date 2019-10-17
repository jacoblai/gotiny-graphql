package engine

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"models"
	"time"
)

func (d *DbEngine) Search(ctx context.Context, args struct{ Name string }) ([]*models.Person, error) {
	c := d.GetColl(models.T_Person)

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
	c := d.GetColl(models.T_Person)

	p := models.Person{
		Name:           args.Input.Name,
		Role:           args.Input.Role,
		Address:        args.Input.Address,
		Email:          args.Input.Email,
		Phone:          args.Input.Phone,
		Total:          args.Input.Total,
		Order:          args.Input.Order,
		CreatedAtFiled: time.Now().Local(),
	}

	id := ""
	dbRes, err := c.InsertOne(context.Background(), &p)
	if err != nil {
		return &id, err
	}

	id = dbRes.InsertedID.(primitive.ObjectID).Hex()

	return &id, nil
}
