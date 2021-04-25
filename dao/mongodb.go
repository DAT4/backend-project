package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/DAT4/backend-project/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	uri string
}

type findOneQuery struct {
	model      interface{}
	filter     bson.M
	options    *options.FindOneOptions
	collection string
}

type addOneQuery struct {
	model      interface{}
	filter     bson.M
	collection string
}

func connect(col string) (*mongo.Collection, *mongo.Client, error) {
	opt := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(opt)
	if err != nil {
		return nil, nil, err
	}
	err = client.Connect(context.Background())
	if err != nil {
		client.Disconnect(context.Background())
		return nil, nil, err
	}
	collection := client.Database("backend").Collection(col)
	return collection, client, nil
}

func (query *findOneQuery) find() error {
	col, cli, err := connect(query.collection)
	if err != nil {
		return err
	}
	defer cli.Disconnect(context.Background())
	if query.options == nil {
		err = col.FindOne(context.Background(), query.filter, options.FindOne()).Decode(query.model)
	} else {
		err = col.FindOne(context.Background(), query.filter, query.options).Decode(query.model)
	}
	if err != nil {
		return err
	}
	return nil
}

func (query *addOneQuery) add() error {
	col, cli, err := connect(query.collection)
	defer cli.Disconnect(context.Background())
	if err != nil {
		return err
	}
	_, err = col.InsertOne(context.Background(), query.model)
	return err
}

func (m *MongoDB) Create(u *models.User) (err error) {
	q2 := addOneQuery{
		model:      &u,
		filter:     nil,
		collection: "users",
	}
	return q2.add()
}
func (m *MongoDB) UserFromId(id string) (user models.User, err error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	q := findOneQuery{
		model:      &user,
		filter:     bson.M{"_id": _id},
		collection: "users",
	}
	err = q.find()
	fmt.Println(user)
	return
}
func (m *MongoDB) Authenticate(u *models.User) error {
	var tmpUser models.User
	q := findOneQuery{
		model: &tmpUser,
		filter: bson.M{
			"username": u.Username,
		},
		collection: "users",
	}

	err := q.find()
	if err != nil {
		return err
	}

	//TODO This logic should be in the middle package
	ok := u.Check(tmpUser.Password)
	if !ok {
		return errors.New("password incorrect")
	}
	u.Id = tmpUser.Id
	u.PlayerID = tmpUser.PlayerID
	u.Password = tmpUser.Password
	return nil
}
func (m *MongoDB) UsernameTaken(u *models.User) (err error) {
	var tmpUser models.User
	q1 := findOneQuery{
		model:      &tmpUser,
		filter:     bson.M{"username": u.Username},
		options:    options.FindOne(),
		collection: "users",
	}
	err = q1.find()
	if err == nil {
		return errors.New("A user already exists with this name")
	}
	return nil
}
