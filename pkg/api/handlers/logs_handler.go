package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Ricky004/watchdata/pkg/clickhousestore"
)


var provider *clickhousestore.ClickHouseProvider

func GetLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logs, err := provider.GetTop10Logs(ctx)
	if err != nil {
		log.Printf("GetTop10Logs error: %v\n", err)
		http.Error(w, "Failed to fetch logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(logs); err != nil {
		http.Error(w, "Failed to encode logs", http.StatusInternalServerError)
	}
}
