// main.go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// -----------------------------
// Main
// -----------------------------

func main() {
	// Load environment variables from .env if it exists
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	port = "8090"

	// Register handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/auth", authHandler)
	mux.HandleFunc("/auth/users/", getUserHandler)

	// Wrap with logging middleware
	loggedMux := loggingMiddleware(mux)

	log.Printf("Remote Authenticator listening on port %s\n", port)
	if err := http.ListenAndServe(":"+port, loggedMux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
