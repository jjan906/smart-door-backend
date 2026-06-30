package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client *mongo.Client
	DB     *mongo.Database
)

// InitMongoDB menginisialisasi koneksi ke MongoDB Atlas / lokal
func InitMongoDB(uri string, dbName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("gagal konek ke MongoDB: %v", err)
	}

	// Ping untuk memastikan koneksi berhasil
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("gagal ping MongoDB: %v", err)
	}

	Client = client
	DB = client.Database(dbName)

	fmt.Println("[MongoDB] Connected to database:", dbName)
	return nil
}

// GetCollection helper untuk ambil collection
func GetCollection(name string) *mongo.Collection {
	return DB.Collection(name)
}
