package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Ricky004/watchdata/pkg/clickhousestore"
)

type Server struct {
    provider *clickhousestore.ClickHouseProvider
}

func NewServer(cfg clickhousestore.Config) (*Server, error) {
    provider, err := clickhousestore.NewClickHouseProvider(context.Background(), cfg)
    if err != nil {
        return nil, fmt.Errorf("failed to create provider: %w", err)
    }
    
    return &Server{provider: provider}, nil
}

func (s *Server) GetLogs(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()

	logs, err := s.provider.GetTop10Logs(ctx)
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
