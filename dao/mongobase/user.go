package mongobase

import (
	"context"
	"github.com/DAT4/backend-project/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Users struct {
	mongo *Mongo
	col   string
}

func NewUsers(mongo *Mongo) *Users {
	return &Users{
		mongo: mongo,
		col:   "users",
	}
}

func (m *Users) Insert(i dto.Object) (err error) {
	col, cli, err := m.mongo.connect(m.col)
	if err != nil {
		return
	}
	defer cli.Disconnect(context.Background())
	out, err := col.InsertOne(context.Background(), i)
	if err != nil {
		return
	}
	err = col.FindOne(context.Background(), bson.M{"_id": out.InsertedID}).Decode(i)
	return
}
func (m *Users) Update(id string, u dto.Update) (o dto.Object, err error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	col, cli, err := m.mongo.connect(m.col)
	if err != nil {
		return
	}
	defer cli.Disconnect(context.Background())
	_, err = col.UpdateOne(context.Background(), bson.M{"_id": _id}, UpdateToBson(u))
	if err != nil {
		return
	}
	err = col.FindOne(context.Background(), bson.M{"_id": _id}).Decode(&o)
	return
}
func (m *Users) Delete(id string) (err error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	col, cli, err := m.mongo.connect(m.col)
	if err != nil {
		return err
	}
	defer cli.Disconnect(context.Background())
	_, err = col.DeleteOne(context.Background(), bson.M{"_id": _id})
	return
}
func (m *Users) FindOne(id string) (o dto.Object, err error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	col, cli, err := m.mongo.connect(m.col)
	if err != nil {
		return
	}
	defer cli.Disconnect(context.Background())
	err = col.FindOne(context.Background(), bson.M{"_id": _id}, options.FindOne()).Decode(&o)
	return
}
func (m *Users) Find(filter dto.Filter) (o []dto.Object, err error) {
	f, err := FilterToBson(filter)
	if err != nil {
		return
	}
	col, cli, err := m.mongo.connect(m.col)
	if err != nil {
		return
	}
	defer cli.Disconnect(context.Background())
	cur, err := col.Find(context.Background(), f, options.Find())
	if err != nil {
		return nil, err
	}
	for cur.TryNext(context.Background()) {
		var user dto.User
		err = cur.Decode(&user)
		if err != nil {
			return
		}
		o = append(o, user)
	}
	return
}
