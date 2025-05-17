package main

import (
	"log"
	"net/http"

	"github.com/Ricky004/watchdata/internals/ingest"
)

func main() {
	http.HandleFunc("/v1/logs", ingest.HandleOTLPLogs)

	log.Println("WatchData OTLP receiver running at :4320/v1/logs")
	if err := http.ListenAndServe(":4320", nil); err != nil {
		log.Fatal(err)
	}
}
