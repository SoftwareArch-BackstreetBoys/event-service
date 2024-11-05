package configs

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure cancel is called to avoid context leak

	// Connect to MongoDB using the mongo.Connect method
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file in server/setup/connectDB")
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGOURI")))
	if err != nil {
		log.Fatal(err)
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")
	return client
}

// Client instance
var DB *mongo.Client = ConnectDB()

// getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("Event").Collection(collectionName)
	return collection
}
