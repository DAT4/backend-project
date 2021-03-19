package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type FindOneQuery struct {
	Model      interface{}
	Filter     bson.M
	Options    *options.FindOneOptions
	Collection string
}

type AddOneQuery struct {
	Model      interface{}
	Filter     bson.M
	Collection string
}

func connect(col string) (*mongo.Collection, *mongo.Client, error) {
	client, err := mongo.NewClient()
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

func (query *FindOneQuery) Find() error {
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

func (query *AddOneQuery) Add() error {
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
