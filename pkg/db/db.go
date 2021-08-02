package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoClient returns mongo.Client
var MongoClient *mongo.Database

// StartConnection starts a connection with a mongodb database
func StartConnection(uri string, appName string) error {
	connectionOptions := options.Client().ApplyURI(uri)
	connectionOptions.SetAppName(appName).SetConnectTimeout(10 * time.Second)
	connectionOptions.SetMaxConnIdleTime(15 * time.Second)
	connectionOptions.SetServerSelectionTimeout(10 * time.Second)

	client, err := mongo.NewClient(connectionOptions)
	if err != nil {
		return err
	}

	err = client.Connect(context.Background())

	if err != nil {
		return err
	}

	err = client.Ping(context.Background(), readpref.Primary())

	if err != nil {
		return err
	}

	MongoClient = client.Database("lkp_prod")

	return nil
}
