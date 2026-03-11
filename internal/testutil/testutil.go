package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// MongoDBConfig holds MongoDB test configuration
type MongoDBConfig struct {
	URI      string
	Database string
	Timeout  time.Duration
}

// DefaultMongoDBConfig returns default MongoDB test configuration
func DefaultMongoDBConfig() MongoDBConfig {
	uri := os.Getenv("TEST_MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	return MongoDBConfig{
		URI:      uri,
		Database: "gohome_test",
		Timeout:  30 * time.Second,
	}
}

// TestDB holds test database connection
type TestDB struct {
	Client   *mongo.Client
	Database *mongo.Database
	DBName   string
	t        *testing.T
}

// NewTestDB creates a new test database connection
func NewTestDB(t *testing.T) *TestDB {
	return NewTestDBWithConfig(t, DefaultMongoDBConfig())
}

// NewTestDBWithConfig creates a new test database connection with custom config
func NewTestDBWithConfig(t *testing.T, config MongoDBConfig) *TestDB {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.URI)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		t.Fatalf("Failed to ping MongoDB: %v", err)
	}

	dbName := config.Database
	if dbName == "" {
		dbName = "gohome_test_" + bson.NewObjectID().Hex()
	}

	db := client.Database(dbName)

	t.Cleanup(func() {
		// Drop the test database
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := db.Drop(ctx); err != nil {
			t.Logf("Warning: failed to drop test database: %v", err)
		}
		if err := client.Disconnect(ctx); err != nil {
			t.Logf("Warning: failed to disconnect from MongoDB: %v", err)
		}
	})

	return &TestDB{
		Client:   client,
		Database: db,
		DBName:   dbName,
		t:        t,
	}
}

// Collection returns a collection from the test database
func (tdb *TestDB) Collection(name string) *mongo.Collection {
	return tdb.Database.Collection(name)
}

// CleanupCollections cleans up specified collections after each test
func (tdb *TestDB) CleanupCollections(collections ...string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, name := range collections {
		if err := tdb.Database.Collection(name).Drop(ctx); err != nil {
			// Ignore "namespace not found" error
			tdb.t.Logf("Warning: failed to drop collection %s: %v", name, err)
		}
	}
}

// SkipIfShort skips the test if -short flag is provided
func SkipIfShort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
}

// SkipIfNoMongoDB skips the test if MongoDB is not available
func SkipIfNoMongoDB(t *testing.T) {
	SkipIfShort(t)

	uri := os.Getenv("TEST_MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		t.Skipf("MongoDB not available: %v", err)
	}
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, nil); err != nil {
		t.Skipf("MongoDB not available: %v", err)
	}
}

// GenerateTestID generates a test ID
func GenerateTestID() string {
	return bson.NewObjectID().Hex()
}

// AssertEqual asserts two values are equal
func AssertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

// AssertNotNil asserts value is not nil
func AssertNotNil(t *testing.T, value interface{}) {
	t.Helper()
	if value == nil {
		t.Error("Expected non-nil value")
	}
}

// AssertError asserts error occurred
func AssertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// AssertNoError asserts no error occurred
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

// WaitForCondition waits for a condition to be true or times out
func WaitForCondition(t *testing.T, timeout time.Duration, condition func() bool, message string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			t.Fatalf("Timeout waiting for condition: %s", message)
		case <-ticker.C:
			if condition() {
				return
			}
		}
	}
}

// RedisConfig holds Redis test configuration
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// DefaultRedisConfig returns default Redis test configuration
func DefaultRedisConfig() RedisConfig {
	addr := os.Getenv("TEST_REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	return RedisConfig{
		Addr:     addr,
		Password: "",
		DB:       15, // Use DB 15 for testing
	}
}

// TestRedisConfig returns Redis config for testing
func TestRedisConfig() string {
	addr := os.Getenv("TEST_REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	return fmt.Sprintf("redis://%s/15", addr)
}