package engine

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"models"
	"time"
)

func (d *DbEngine) SearchOrders(ctx context.Context, args struct {
	Name        string
	Skip, Limit *int32
}) ([]*models.Person, error) {
	c := d.GetColl(models.T_Person)

	qstr := bson.M{"name": args.Name}
	query := []bson.M{
		{"$match": qstr},
		{"$sort": bson.M{"createdat": -1}},
	}

	if args.Skip != nil {
		query = append(query, bson.M{"$skip": args.Skip})
	}
	if args.Limit != nil {
		query = append(query, bson.M{"$limit": *args.Limit})
	}
	query = append(query,
		bson.M{"$lookup": bson.M{
			"from":         models.T_Order,
			"localField":   "_id",
			"foreignField": "personid",
			"as":           "orders",
		}})

	var objs []*models.Person
	re, err := c.Aggregate(ctx, query, options.Aggregate())
	if err != nil {
		return nil, err
	}
	err = re.All(ctx, &objs)
	if err != nil {
		return nil, err
	}

	return objs, nil
}

func (d *DbEngine) CreateOrder(ctx context.Context, args struct {
	PersonId string
	Input    *models.InputOrder
}) (*string, error) {
	c := d.GetColl(models.T_Order)

	oid, err := primitive.ObjectIDFromHex(args.PersonId)
	if err != nil {
		return nil, err
	}

	o := models.Order{
		PersonIdFiled:  oid,
		Express:        args.Input.Express,
		IsDisable:      args.Input.IsDisable,
		CreatedAtFiled: time.Now().Local(),
	}

	id := ""
	dbRes, err := c.InsertOne(context.Background(), &o)
	if err != nil {
		return &id, err
	}

	id = dbRes.InsertedID.(primitive.ObjectID).Hex()

	return &id, nil
}
