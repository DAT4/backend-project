package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
	opt := options.Client().ApplyURI("mongodb://mongo:27017")
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
