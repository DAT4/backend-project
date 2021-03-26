package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type findOneQuery struct {
	Model      interface{}
	Filter     bson.M
	Options    *options.FindOneOptions
	Collection string
}

type addOneQuery struct {
	Model      interface{}
	Filter     bson.M
	Collection string
}

func connect(col string) (*mongo.Collection, *mongo.Client, error) {
	client, err := mongo.NewClient()
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
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	col, cli, err := connect(query.Collection)
	if err != nil {
		return err
	}
	defer cli.Disconnect(ctx)
	if query.Options == nil {
		err = col.FindOne(ctx, query.Filter, options.FindOne()).Decode(query.Model)
	} else {
		err = col.FindOne(ctx, query.Filter, query.Options).Decode(query.Model)
	}
	if err != nil {
		return err
	}
	return nil
}

func (query *addOneQuery) add() error {
	col, cli, err := connect(query.Collection)
	defer cli.Disconnect(context.Background())
	if err != nil {
		return err
	}
	_, err = col.InsertOne(context.Background(), query.Model)
	return err
}
