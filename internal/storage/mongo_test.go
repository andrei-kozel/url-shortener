package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/andrei-kozel/url-shortener/internal/model"
	"github.com/andrei-kozel/url-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupMongoDB() *mongo.Client {
	// Connect to MongoDB for testing purposes
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	return client
}

func teardownMongoDB(client *mongo.Client) {
	// Disconnect from MongoDB
	err := client.Disconnect(context.Background())
	if err != nil {
		panic(err)
	}
}

func TestMongoDBStorage_PutAndGet(t *testing.T) {
	client := setupMongoDB()
	defer teardownMongoDB(client)

	storage := storage.NewMongoDB(client)

	// Test data
	shortening := model.Shortening{
		Identifier:  "abc123",
		OriginalUrl: "http://example.com",
		Visits:      0,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// check if the identifier already exists
	getShortening, err := storage.Get(context.Background(), shortening.Identifier)
	if err != nil {
		panic(err)
	}
	if getShortening != nil {
		// Delete the existing shortening
		err := storage.Delete(context.Background(), shortening.Identifier)
		assert.NoError(t, err)
	}

	// Test Put
	storedShortening, err := storage.Put(context.Background(), shortening)
	assert.NoError(t, err)
	assert.NotNil(t, storedShortening)

	// Test Get
	retrievedShortening, err := storage.Get(context.Background(), shortening.Identifier)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedShortening)
	assert.Equal(t, shortening.Identifier, retrievedShortening.Identifier)
	assert.Equal(t, shortening.OriginalUrl, retrievedShortening.OriginalUrl)
	assert.Equal(t, shortening.Visits, retrievedShortening.Visits)
	assert.Equal(t, shortening.CreatedAt.Format("2006-01-01"), retrievedShortening.CreatedAt.Format("2006-01-02"))
	assert.Equal(t, shortening.UpdatedAt.Format("2006-01-01"), retrievedShortening.UpdatedAt.Format("2006-01-02"))
}

func TestMongoDBStorage_IncrementVisits(t *testing.T) {
	client := setupMongoDB()
	defer teardownMongoDB(client)

	storage := storage.NewMongoDB(client)

	// Test data
	identifier := "abc123"

	// Test IncrementVisits
	err := storage.IncrementVisits(context.Background(), identifier)
	assert.NoError(t, err)
}
