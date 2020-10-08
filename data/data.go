package data

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetData . . .
func GetData() (*mongo.Client, error) {

	clientOptions := options.Client().ApplyURI(os.Getenv("ConnString"))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client, err

}
