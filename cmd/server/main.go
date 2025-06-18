package main

import (
	"log"
	"net/http"

	"github.com/Ricky004/watchdata/pkg/api/handlers"
	"github.com/Ricky004/watchdata/pkg/clickhousestore"
)

func main() {
	// Load ClickHouse configuration
	cfg, err := clickhousestore.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize server with ClickHouse provider
	server, err := handlers.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Register routes
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/logs", server.GetLogs)
	mux.HandleFunc("/v1/logs/since", server.GetLogsSince)
	mux.HandleFunc("/v1/logs/timerange", server.GetLogsInTimeRanges)
	mux.HandleFunc("/ws", server.WebSocketHandler)

	log.Println("🚀 Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
