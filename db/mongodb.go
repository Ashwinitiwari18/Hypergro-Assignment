package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client          *mongo.Client
	Database        *mongo.Database
	Properties      *mongo.Collection
	Users           *mongo.Collection
	Favorites       *mongo.Collection
	Recommendations *mongo.Collection
)

func InitDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connectionString := os.Getenv("MONGODB_URI")
	if connectionString == "" {
		return fmt.Errorf("MONGODB_URI environment variable is not set")
	}

	log.Printf("Attempting to connect to MongoDB...")

	clientOptions := options.Client().
		ApplyURI(connectionString).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)).
		SetConnectTimeout(5 * time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("Error connecting to MongoDB: %v", err)
		return fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Error pinging MongoDB: %v", err)
		return fmt.Errorf("error pinging MongoDB: %v", err)
	}

	log.Printf("Successfully connected to MongoDB!")

	Client = client
	Database = client.Database("property_listing")
	Properties = Database.Collection("properties")
	Users = Database.Collection("users")
	Favorites = Database.Collection("favorites")
	Recommendations = Database.Collection("recommendations")

	// Create indexes
	if err := createIndexes(); err != nil {
		log.Printf("Warning: Error creating indexes: %v", err)
	}

	return nil
}

func createIndexes() error {
	ctx := context.Background()

	// Users collection indexes
	_, err := Users.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    map[string]interface{}{"email": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		log.Printf("Error creating email index: %v", err)
		return err
	}

	// Properties collection indexes
	_, err = Properties.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{"createdBy": 1},
	})
	if err != nil {
		log.Printf("Error creating createdBy index: %v", err)
		return err
	}

	// Favorites collection indexes
	_, err = Favorites.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{
			"userId":     1,
			"propertyId": 1,
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		log.Printf("Error creating favorites index: %v", err)
		return err
	}

	// Recommendations collection indexes
	_, err = Recommendations.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{
			"toUserId": 1,
			"isRead":   1,
		},
	})
	if err != nil {
		log.Printf("Error creating recommendations index: %v", err)
		return err
	}

	return nil
}
