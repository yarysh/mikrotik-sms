package main

import (
	"log/slog"
	"net/http"
	"os"
)

func main() {
	port := getEnvDefault("PORT", "8080")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK!\n"))
	})

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
