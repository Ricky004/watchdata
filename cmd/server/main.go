package main

import (
	"log"
	"net/http"

	"github.com/Ricky004/watchdata/pkg/api/routes"
)

func main() {

	log.Println("🚀 Server running at :8080")
	http.ListenAndServe(":8080", routes.RegisterRoutes())
	
}
