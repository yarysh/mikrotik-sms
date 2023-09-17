package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	apiKey := getEnvDefault("API_KEY", "secret")
	mongoDBURI := getEnvDefault("MONGODB_URI", "mongodb://localhost:27017")
	port := getEnvDefault("PORT", "8080")

	mongoClient, err := connectDB(mongoDBURI)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			slog.Error("error disconnecting from MongoDB: ", err)
		}
	}()

	api := NewAPI(apiKey, mongoClient.Database("sms"))
	http.HandleFunc("/incoming_sms", api.IncomingSMS)

	slog.Info("Starting server on port " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		slog.Error("ListenAndServe: ", err)
		os.Exit(1)
	}
}

func getEnvDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func connectDB(uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("error creating MongoDB client: %w", err)
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	return client, nil
}
