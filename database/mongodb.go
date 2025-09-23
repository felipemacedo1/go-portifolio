package database

import (
	"context"
	"log"
	"portfolio-backend/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client   *mongo.Client
	Database *mongo.Database
)

func Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create MongoDB client
	clientOptions := options.Client().ApplyURI(config.AppConfig.MongoDBURI)
	
	var err error
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Test the connection
	err = Client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	// Get database instance
	Database = Client.Database(config.AppConfig.DatabaseName)

	log.Printf("Connected to MongoDB: %s", config.AppConfig.DatabaseName)
	
	// Create indexes
	if err := createIndexes(); err != nil {
		log.Printf("Warning: Failed to create indexes: %v", err)
	}

	return nil
}

func Disconnect() error {
	if Client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return Client.Disconnect(ctx)
}

func createIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create TTL index for cache collection
	cacheCollection := Database.Collection("cache")
	indexModel := mongo.IndexModel{
		Keys:    map[string]interface{}{"expires_at": 1},
		Options: options.Index().SetExpireAfterSeconds(0),
	}
	
	_, err := cacheCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}

	// Create index for content collection
	contentCollection := Database.Collection("content")
	contentIndexModel := mongo.IndexModel{
		Keys: map[string]interface{}{"type": 1, "updated_at": -1},
	}
	
	_, err = contentCollection.Indexes().CreateOne(ctx, contentIndexModel)
	if err != nil {
		return err
	}

	log.Println("Database indexes created successfully")
	return nil
}

// Health check function
func IsHealthy() bool {
	if Client == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := Client.Ping(ctx, nil)
	return err == nil
}