package src

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

var db *mongo.Database

const dbSource = "mongodb://localhost:27017"
const dbName = "golf"

func init() {
	fmt.Println("conn info:", dbSource, dbName)

	var logger log.Logger
	clientOptions := options.Client()
	clientOptions.ApplyURI(dbSource)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		level.Error(logger).Log("exit", err)
		os.Exit(-1)
	}

	err = client.Connect(context.Background())
	if err != nil {
		level.Error(logger).Log("exit", err)
		os.Exit(-1)
	}

	db = client.Database(dbName)
}

// GetMongoDB function to return DB connection
func GetMongoDB() *mongo.Database {
	return db
}
