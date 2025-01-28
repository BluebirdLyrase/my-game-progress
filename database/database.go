package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Context context.Context
var DB *mongo.Database

func Init() error {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Create a MongoDB client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		cancel() // Clean up context if connection fails
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Verify the connection
	if err := client.Ping(ctx, nil); err != nil {
		cancel()
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Connected to MongoDB")
	Client = client
	DB = client.Database("my-game-progress")
	Context = ctx
	// collection := database.Collection("Game")
	return nil
}

func Fetch(collection *mongo.Collection) ([]bson.M, error) {
	cursor, err := collection.Find(Context, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(Context)

	var table []bson.M
	for cursor.Next(Context) {
		var row bson.M
		if err := cursor.Decode(&row); err != nil {
			return nil, fmt.Errorf("failed to decode document: %w", err)
		}
		table = append(table, row)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return table, nil
}
