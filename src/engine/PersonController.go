package engine

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	id := ""

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

	err := d.MgEngine.UseSessionWithOptions(context.Background(), options.Session().SetDefaultReadPreference(readpref.Primary()), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}

		dbRes, err := c.InsertOne(sessionContext, &p)
		if err != nil {
			_ = sessionContext.AbortTransaction(sessionContext)
			return err
		}

		id = dbRes.InsertedID.(primitive.ObjectID).Hex()

		return sessionContext.CommitTransaction(sessionContext)
	})
	if err != nil {
		return nil, err
	}

	return &id, nil
}
