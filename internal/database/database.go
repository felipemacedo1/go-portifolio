package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB represents the database connection
type DB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// Connect establishes a connection to MongoDB
func Connect(uri string) (*DB, error) {
	// Don't attempt connection to localhost in development
	if uri == "mongodb://localhost:27017" {
		return nil, fmt.Errorf("skipping localhost connection in development mode")
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	database := client.Database("portfolio")

	return &DB{
		Client:   client,
		Database: database,
	}, nil
}

// Disconnect closes the database connection
func Disconnect(db *DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return db.Client.Disconnect(ctx)
}

// GetCollection returns a collection from the database
func (db *DB) GetCollection(name string) *mongo.Collection {
	return db.Database.Collection(name)
}