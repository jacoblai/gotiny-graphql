package engine

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"models"
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
		{"$skip": args.Skip},
	}

	if *args.Limit > 0 {
		query = append(query, bson.M{"$limit": *args.Limit})
	}
	query = append(query,
		bson.M{"$lookup": bson.M{
			"from":         models.T_Order,
			"localField":   "personid",
			"foreignField": "_id",
			"as":           "order",
		}})

	var objs []*models.Person
	re, err := c.Aggregate(ctx, query)
	if err != nil {
		return nil, err
	}
	err = re.All(ctx, &objs)
	if err != nil {
		return nil, err
	}

	return objs, nil
}
