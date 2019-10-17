package models

import (
	"github.com/graph-gophers/graphql-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var T_Order = "orders"

//订单
type Order struct {
	IdFiled        primitive.ObjectID `json:"id" bson:"_id,omitempty" jsonschema:"-"`
	CreatedAtFiled time.Time          `json:"createdat,omitempty" bson:"createdat,omitempty"`

	//用户信息
	PersonIdFiled primitive.ObjectID `json:"personid,omitempty" bson:"personid,omitempty" jsonschema:"required,oid"`
	Express       string             `json:"express,omitempty" bson:"express,omitempty" jsonschema:"required,enum=Wait|WaitDelivery|Shipped|Finish|Cancel"`
	IsDisable     bool               `json:"isdisable,omitempty" bson:"isdisable,omitempty"` //是否冻结
}

func (u Order) CreatedAt() graphql.Time {
	return graphql.Time{Time: u.CreatedAtFiled.Local()}
}

func (u Order) Id() string {
	return u.IdFiled.Hex()
}

func (u Order) PersonId() string {
	return u.PersonIdFiled.Hex()
}
