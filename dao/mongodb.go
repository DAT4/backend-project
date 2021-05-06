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

func (q *query) addOne() (id primitive.ObjectID, err error) {
	col, cli, err := q.db.connect(q.collection)
	if err != nil {
		return
	}
	defer cli.Disconnect(context.Background())
	x, err := col.InsertOne(context.Background(), q.model)
	if err != nil {
		return
	}
	id, ok := x.InsertedID.(primitive.ObjectID)
	if ok {
		return
	}
	return id, errors.New("mongo id could not compile to primitive")
}

func (m *MongoDB) Create(u *models.User) (err error) {
	q2 := query{
		db:         m,
		model:      &u,
		filter:     nil,
		collection: "users",
	}
	id, err := q2.addOne()
	if err != nil {
		return err
	}
	col, cli, err := m.connect("users")
	if err != nil {
		return err
	}
	defer cli.Disconnect(context.Background())
	newId, err := m.getNextSequence("users")
	if err != nil {
		return
	}
	x := col.FindOneAndUpdate(context.Background(), bson.M{"_id": id}, bson.M{"$set": bson.M{"playerid": newId.Seq}})
	err = x.Decode(u)
	return
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

type counter struct {
	Id  string `bson:"_id"`
	Seq int    `bson:"seq"`
}

func (m *MongoDB) getNextSequence(name string) (out counter, err error) {
	q := bson.M{"_id": name}
	u := bson.M{"$inc": bson.M{"seq": 1}}
	col, cli, err := m.connect("counters")
	if err != nil {
		return
	}
	defer cli.Disconnect(context.Background())
	err = col.FindOneAndUpdate(context.Background(), q, u).Decode(&out)
	return
}
