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
	Uri string
}

type _query struct {
	model      interface{}
	filter     bson.M
	options    *options.FindOneOptions
	collection string
}

type query struct {
	model      interface{}
	db         *MongoDB
	filter     bson.M
	collection string
}

func (m *MongoDB) connect(col string) (*mongo.Collection, *mongo.Client, error) {
	fmt.Println(m.Uri)
	opt := options.Client().ApplyURI(m.Uri)
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

func (q *query) findOne(o *options.FindOneOptions) error {
	if o == nil {
		o = options.FindOne()
	}
	col, cli, err := q.db.connect(q.collection)
	if err != nil {
		return err
	}
	defer cli.Disconnect(context.Background())
	err = col.FindOne(context.Background(), q.filter, o).Decode(q.model)
	if err != nil {
		return err
	}
	return nil
}

func (q *query) addOne() error {
	col, cli, err := q.db.connect(q.collection)
	defer cli.Disconnect(context.Background())
	if err != nil {
		return err
	}
	_, err = col.InsertOne(context.Background(), q.model)
	return err
}

func (m *MongoDB) Create(u *models.User) (err error) {
	q2 := query{
		model:      &u,
		filter:     nil,
		collection: "users",
	}
	return q2.addOne()
}
func (m *MongoDB) UserFromId(id string) (user models.User, err error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	q := query{
		db:         m,
		model:      &user,
		filter:     bson.M{"_id": _id},
		collection: "users",
	}
	err = q.findOne(nil)
	fmt.Println(user)
	return
}
func (m *MongoDB) Authenticate(u *models.User) error {
	var tmpUser models.User
	q := query{
		db:    m,
		model: &tmpUser,
		filter: bson.M{
			"username": u.Username,
		},
		collection: "users",
	}

	err := q.findOne(nil)
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
	q1 := query{
		db:         m,
		model:      &tmpUser,
		filter:     bson.M{"username": u.Username},
		collection: "users",
	}
	err = q1.findOne(nil)
	if err == nil {
		return errors.New("A user already exists with this name")
	}
	return nil
}
