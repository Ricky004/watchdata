package ingest

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"

	collectorpb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	"google.golang.org/protobuf/proto"
)

func HandleOTLPLogs(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-protobuf") {
		http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	var reader io.Reader = bytes.NewReader(body)
	if r.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(reader)
		if err != nil {
			log.Printf("gzip.NewReader failed: %v", err)
			http.Error(w, "Error decompressing gzip data", http.StatusBadRequest)
			return
		}
		defer gzipReader.Close()
		reader = gzipReader 
		log.Println("Received gzipped data, decompressing.")
	} else {
		log.Println("Received uncompressed data.")
	}

	// Read the (possibly decompressed) body
	decompressedBody, err := io.ReadAll(reader)
	if err != nil {
		log.Printf("io.ReadAll after decompression failed: %v", err)
		http.Error(w, "Failed to read decompressed body", http.StatusBadRequest)
		return
	}

	req := &collectorpb.ExportLogsServiceRequest{}
	err = proto.Unmarshal(decompressedBody, req) // Use decompressedBody
	if err != nil {
		log.Printf("Unmarshal error: %v", err)
		http.Error(w, "Failed to decode OTLP protobuf", http.StatusBadRequest)
		return
	}

	count := 0
	for _, rl := range req.ResourceLogs {
		for _, sl := range rl.ScopeLogs {
			count += len(sl.LogRecords)
		}
	}

	log.Println("Received OTLP logs:", count)

	w.WriteHeader(http.StatusAccepted)
}
