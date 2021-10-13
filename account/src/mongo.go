package src

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoDB(ctx context.Context, dbsource string, dbname string) (*mongo.Database, error) {
	clientOptions := options.Client()
	clientOptions.ApplyURI(dbsource)
	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client.Database(dbname), nil
}
