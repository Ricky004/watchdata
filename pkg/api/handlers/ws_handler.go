package handlers

import (
	"log"
	"net/http"

	"github.com/Ricky004/watchdata/pkg/types/telemetrytypes"
	"github.com/gorilla/websocket"
)

var Clients = make(map[*websocket.Conn]bool)
var Broadcast = make(chan telemetrytypes.LogRecord) 

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Upgrade error:", err)
        return
    }
    Clients[ws] = true

    go func() {
        for {
            if _, _, err := ws.ReadMessage(); err != nil {
                ws.Close()
                delete(Clients, ws)
                break
            }
        }
    }()
}

func StartBroadcaster() {
    go func() {
        for {
            logRecord := <-Broadcast
            for client := range Clients {
                err := client.WriteJSON(logRecord)
                if err != nil {
                    client.Close()
                    delete(Clients, client)
                }
            }
        }
    }()
}
