package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Ricky004/watchdata/pkg/clickhousestore"
	"github.com/Ricky004/watchdata/pkg/types/telemetrytypes"
	"github.com/gorilla/websocket"
)

type Server struct {
	provider  *clickhousestore.ClickHouseProvider
	clients   map[*websocket.Conn]bool
	clientsMu sync.Mutex
	broadcast chan telemetrytypes.LogRecord
	upgrader  websocket.Upgrader
}

func NewServer(cfg clickhousestore.Config) (*Server, error) {
	provider, err := clickhousestore.NewClickHouseProvider(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}

	server := &Server{
		provider:  provider,
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan telemetrytypes.LogRecord, 1000), // Buffer for high throughput
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}

	// Start the broadcaster goroutine
	server.startBroadcaster()
	server.startDatabasePoller()

	return server, nil
}

func (s *Server) startDatabasePoller() {
	go func() {
		log.Println("Starting database poller...")

		// Use current time minus a bit for initial buffer
		lastTimestamp := time.Now().Add(-5 * time.Second)

		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			newLogs, err := s.provider.GetLogsSince(context.Background(), lastTimestamp)
			if err != nil {
				log.Printf("Error polling for new logs: %v", err)
				continue
			}

			if len(newLogs) == 0 {
				continue
			}

			filteredLogs := []telemetrytypes.LogRecord{}

			for _, logRecord := range newLogs {
				// skip logs that have the same timestamp as last seen
				if logRecord.Timestamp.After(lastTimestamp) {
					filteredLogs = append(filteredLogs, logRecord)
				}
			}

			if len(filteredLogs) > 0 {
				log.Printf("Found %d new logs to broadcast.", len(filteredLogs))
				for _, logRecord := range filteredLogs {
					s.broadcast <- logRecord
				}

				// update to the highest timestamp seen
				lastTimestamp = filteredLogs[len(filteredLogs)-1].Timestamp
			}
		}
	}()
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

// WebSocket handler
func (s *Server) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w)

	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	// Add client to the map
	s.clientsMu.Lock()
	s.clients[ws] = true
	log.Printf("WebSocket client connected. Total clients: %d", len(s.clients))
	s.clientsMu.Unlock()

	// Handle client disconnection
	go func() {
		defer func() {
			s.clientsMu.Lock()
			delete(s.clients, ws)
			log.Printf("WebSocket client disconnected. Total clients: %d", len(s.clients))
			s.clientsMu.Unlock()
			ws.Close()
		}()

		// Keep connection alive and handle client messages (if any)
		for {
			// Read message from client (ping/pong or other messages)
			if _, _, err := ws.ReadMessage(); err != nil {
				log.Printf("WebSocket read error: %v", err)
				break
			}
		}
	}()
}

// Start the broadcaster goroutine
func (s *Server) startBroadcaster() {
	go func() {
		log.Println("WebSocket broadcaster started")
		for logRecord := range s.broadcast {
			// FIX: Lock the mutex to safely read the clients map
			s.clientsMu.Lock()

			if len(s.clients) == 0 {
				s.clientsMu.Unlock()
				continue // No clients connected, skip broadcasting
			}

			log.Printf("Broadcasting log to %d clients", len(s.clients))

			// FIX: To prevent issues with modifying the map while iterating,
			// we collect clients that have disconnected and remove them after the loop.
			var badClients []*websocket.Conn

			// Send to all connected clients
			for client := range s.clients {
				err := client.WriteJSON(logRecord)
				if err != nil {
					log.Printf("Error sending to WebSocket client: %v", err)
					badClients = append(badClients, client)
				}
			}

			// Remove any clients that failed to receive the message
			for _, client := range badClients {
				client.Close()
				delete(s.clients, client)
			}

			s.clientsMu.Unlock()
		}
	}()
}

// IngestLog - called by your collector/exporter (send logs directly to the server's API)
func (s *Server) IngestLog(logs []telemetrytypes.LogRecord) {
	// Store to ClickHouse first
	err := s.provider.InsertLogs(context.Background(), logs)
	if err != nil {
		log.Printf("Error storing logs to ClickHouse: %v", err)
		return // Don't broadcast if storage failed
	}

	// Broadcast to WebSocket clients
	for _, logRecord := range logs {
		select {
		case s.broadcast <- logRecord:
			// Successfully queued for broadcast
		default:
			// This might happen if the channel is full.
			log.Println("Warning: Broadcast channel full, dropping log")
		}
	}
}
