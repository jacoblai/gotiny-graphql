package engine

import (
	"context"
	"github.com/jacoblai/mschema"
	"github.com/pquerna/ffjson/ffjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
	"models"
	"sync"
	"time"
)

type DbEngine struct {
	MgEngine *mongo.Client //关系型数据库引擎
	Mdb      string
	RoleLock sync.RWMutex
}

func NewDbEngine() *DbEngine {
	return &DbEngine{}
}

func (d *DbEngine) Open(dir, mg, mdb string, initdb bool) error {
	d.Mdb = mdb
	ops := options.Client().ApplyURI(mg)
	p := uint64(39000)
	ops.MaxPoolSize = &p
	ops.WriteConcern = writeconcern.New(writeconcern.J(true), writeconcern.W(1))
	ops.ReadPreference = readpref.PrimaryPreferred()
	db, err := mongo.NewClient(ops)
	//db, err := mgo.Dial("")
	//mongodb://root:1q2w3e4r@192.168.100.251:27017,192.168.100.252:27017,192.168.100.250:27017/?authSource=admin&readPreference=primaryPreferred&replicaSet=rs1
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = db.Connect(ctx)
	if err != nil {
		return err
	}
	err = db.Ping(ctx, readpref.PrimaryPreferred())
	if err != nil {
		log.Println("ping err", err)
	}

	d.MgEngine = db

	//初始化数据库
	if initdb {
		var session *mongo.Client
		ss, err := mongo.NewClient(ops)
		if err != nil {
			panic(err)
		}
		err = ss.Connect(context.Background())
		if err != nil {
			panic(err)
		}
		session = ss
		defer session.Disconnect(context.Background())
		//user表
		res := InitDbAndColl(session, mdb, models.T_Person, GenJsonSchema(&models.Person{}))
		u := session.Database(mdb).Collection(models.T_Person)
		indexView := u.Indexes()
		_, err = indexView.CreateOne(context.Background(),
			mongo.IndexModel{
				Keys: bsonx.Doc{{"name", bsonx.Int32(1)}},
			})
		if err != nil {
			log.Println(err)
		}
		log.Println(models.T_Person, res["ok"])
		//order表
		res = InitDbAndColl(session, mdb, models.T_Order, GenJsonSchema(&models.Order{}))
		log.Println(models.T_Order, res["ok"])
	}

	return nil
}

func (d *DbEngine) GetSess() (mongo.Session, error) {
	session, err := d.MgEngine.StartSession(options.Session().SetDefaultReadPreference(readpref.Primary()))
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (d *DbEngine) GetColl(coll string) *mongo.Collection {
	col, _ := d.MgEngine.Database(d.Mdb).Collection(coll).Clone()
	return col
}

func InitDbAndColl(session *mongo.Client, db, coll string, model map[string]interface{}) map[string]interface{} {
	tn, _ := session.Database(db).ListCollections(context.Background(), bson.M{"name": coll})
	if tn.Next(context.Background()) == false {
		session.Database(db).RunCommand(context.Background(), bson.D{{"create", coll}})
	}
	result := session.Database(db).RunCommand(context.Background(), bson.D{{"collMod", coll}, {"validator", model}})
	var res map[string]interface{}
	err := result.Decode(&res)
	if err != nil {
		log.Println(err)
	}
	return res
}

//创建数据库验证schema结构对象
func GenJsonSchema(obj interface{}) map[string]interface{} {
	flect := &mschema.Reflector{ExpandedStruct: true, RequiredFromJSONSchemaTags: true, AllowAdditionalProperties: true}
	ob := flect.Reflect(obj)
	bts, _ := ffjson.Marshal(&ob)
	var o map[string]interface{}
	_ = ffjson.Unmarshal(bts, &o)
	return bson.M{"$jsonSchema": o}
}

func (d *DbEngine) Close() {
	d.MgEngine.Disconnect(context.Background())
}
