package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Struct used for maintaining HTTP Request Context
type MongoClient struct {
	Context context.Context
}

var client *mongo.Client
var dbName *string

// XXX Create a new DocumentDb client object
// Added 28-31, and chained the call to .SetTLSConfig(tlsConfig)
// TODO: base off configuration of DBType = "DocDb"
func NewMongoClient(ctx context.Context, mongoDBHost string, database string) (MongoDatabase, error) {
	mongoClient := &MongoClient{Context: ctx}
	var err error
	tlsConfig, err := GetCustomTLSConfig(caFilePath)
	if err != nil {
		log.Fatalf("Failed getting TLS configuration: %v", err)
	}
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoDBHost).SetTLSConfig(tlsConfig))
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping cluster: %v", err)
	}

	dbName = &database
	return mongoClient, err
}

// GetContext get mongodb context
func (c *MongoClient) GetContext() (context.Context, error) {
	if client == nil {
		return nil, fmt.Errorf("MongoDB client is nil, be sure you have invoked NewClient() function already!")
	}
	return c.Context, nil
}

// Close Client
func (c *MongoClient) Close() error {
	return client.Disconnect(c.Context)
}

// GetCollection gets database collection
func (c *MongoClient) GetCollection(collectionName string) (*mongo.Collection, error) {
	if client == nil {
		return nil, fmt.Errorf("MongoDB client is nil, be sure you have invoked NewClient() function already!")
	}
	collection := client.Database(*dbName).Collection(collectionName)

	return collection, nil
}

// Create database session
func (c *MongoClient) GetDb() (*mongo.Database, error) {
	if client == nil {
		return nil, fmt.Errorf("MongoDB client is nil, be sure you have invoked NewClient() function already!")
	}
	db := client.Database(*dbName)

	return db, nil
}

// Ping database
func (c *MongoClient) Ping() error {
	if client == nil {
		_ = client.Ping(c.Context, readpref.Primary())
	}

	return nil
}
