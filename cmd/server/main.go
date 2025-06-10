package main

import (
	"log"
	"net/http"

	"github.com/Ricky004/watchdata/pkg/api/handlers"
	"github.com/Ricky004/watchdata/pkg/clickhousestore"
)

func main() {
    cfg, err := clickhousestore.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
    
    server, err := handlers.NewServer(cfg)
    if err != nil {
        log.Fatal("Failed to create server:", err)
    }

	http.HandleFunc("/v1/logs", server.GetLogs)
    log.Fatal(http.ListenAndServe(":8080", nil))
	
}
