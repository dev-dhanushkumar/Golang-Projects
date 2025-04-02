package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoURI       = "mongodb://mongodb:qwe123PLM@localhost:27017"
	databaseName   = "vanigam"
	collectionName = "products"
	timeout        = 10 * time.Second
)

func GetMongoClient() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("Failed to connect MongoDb: %v", err)
		return nil, err
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Failed to ping MongoDb: %v", err)
		return nil, err
	}

	log.Printf("Connected to MongoDb Successfully")
	return client, nil

}

func GetProductsCollection(client *mongo.Client) *mongo.Collection {
	return client.Database(databaseName).Collection(collectionName)
}
