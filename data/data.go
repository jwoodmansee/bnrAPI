package data

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetData . . .
func GetData() (*mongo.Client, error) {

	clientOptions := options.Client().ApplyURI("mongodb+srv://jwoodmansee:bV81K4XgArtc82x7@cluster0.rsqg9.mongodb.net/Cluster0?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client, err

}
