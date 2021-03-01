package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

func connect(col string) (*mongo.Collection, *mongo.Client, error) {
	uri := os.Getenv("MONGO_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		client.Disconnect(ctx)
		return nil, nil, err
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	collection := client.Database("backend").Collection(col)
	return collection, client, nil
}

func (query *FindOneQuery) find() error {
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

func (query *AddOneQuery) add() error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	col, cli, err := connect(query.Collection)
	defer cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	_, err = col.InsertOne(ctx, query.Model)
	if err != nil {
		return err
	}
	return nil
}
