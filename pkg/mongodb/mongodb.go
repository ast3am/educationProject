package mongodb

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, host, port, database string) (db *mongo.Database, err error) {
	mongoDBURL := fmt.Sprintf("mongodb://%s:%s", host, port)
	clientOptions := options.Client().ApplyURI(mongoDBURL)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		err = errors.New("can't connect")
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		err = errors.New("ping error")
		return nil, err
	}
	return client.Database(database), nil
}
