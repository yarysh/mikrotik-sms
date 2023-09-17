package main

import (
	"io"
	"log/slog"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type API struct {
	apiKey string
	db     *mongo.Database
}

func NewAPI(apiKey string, db *mongo.Database) *API {
	return &API{
		apiKey: apiKey,
		db:     db,
	}
}

func (a *API) IncomingSMS(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed\n"))
		return
	}

	if r.Header.Get("x-api-key") != a.apiKey {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized\n"))
		return
	}

	sms, err := io.ReadAll(io.LimitReader(r.Body, 1024))
	if err != nil {
		slog.Error("error getting incoming SMS: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := a.db.Collection("sms").InsertOne(r.Context(), bson.M{
		"raw":     string(sms),
		"created": time.Now(),
	}); err != nil {
		slog.Error("error saving SMS to DB: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
