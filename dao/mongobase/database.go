package mongobase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DAT4/backend-project/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"strings"
)

type Mongo struct {
	Uri string
}

func (m *Mongo) connect(col string) (*mongo.Collection, *mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(m.Uri)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, nil, err
	}
	err = client.Connect(context.Background())
	if err != nil {
		client.Disconnect(context.Background())
		return nil, nil, err
	}
	collection := client.Database("moro").Collection(col)
	return collection, client, nil
}
func FilterToBson(filter dto.Filter) (out bson.M, err error) {
	if j, err := json.Marshal(filter); err != nil {
		return nil, err
	} else {
		if err = bson.UnmarshalExtJSON(j, true, &out); err != nil {
			return nil, err
		} else {
			fmt.Println(out)
			return out, nil
		}
	}
}
func UpdateToBson(u dto.Update) bson.M {
	out := make(bson.M)
	v := reflect.ValueOf(u)
	t := reflect.TypeOf(u)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).String() != "" {
			out[strings.ToLower(t.Field(i).Name)] = v.Field(i).String()
		}
	}
	return bson.M{"$set": out}
}
