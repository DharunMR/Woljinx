package database

import (
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect() *mongo.Client {

	MongoDb := os.Getenv("MONGODB_URI")
	if MongoDb == "" {
		log.Fatal("MONGODB_URI not set!")
	}
	fmt.Println("MONGODB_URI: ", MongoDb)

	clientOptions := options.Client().ApplyURI(MongoDb)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil
	}

	return client

}

func OpenCollection(collectionName string, client *mongo.Client) *mongo.Collection {

	databaseName := os.Getenv("DATABASE_NAME")

	collection := client.Database(databaseName).Collection(collectionName)

	if collection == nil {
		return nil
	}

	log.Println("Using database:", databaseName)
	log.Println("Using collection:", collection)

	return collection

}
