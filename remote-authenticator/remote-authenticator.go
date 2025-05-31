// main.go
package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// -----------------------------
// Data Types (schemas)
// -----------------------------

// AuthRequestBody represents the request body for POST /auth
type AuthRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponseBody is the response for a successful POST /auth
type AuthResponseBody struct {
	ExternalUserIdentifier string `json:"external_user_identifier"`
}

// ErrorResponse is used for sending error messages
type ErrorResponse struct {
	Message string `json:"message"`
}

// UserDetail is returned by GET /auth/users/{userid}
type UserDetail struct {
	Username          string `json:"username"`
	CreationDate      int64  `json:"creation_date"`
	FriendlyName      string `json:"friendly_name"`
	Email             string `json:"email"`
	Source            string `json:"source"`
	EncryptedPassword string `json:"encryptedPassword"`
	ExternalID        string `json:"external_id"`
}

// -----------------------------
// In-Memory User Store
// -----------------------------
// This map holds UserDetail entries keyed by username.
// A mutex guards concurrent access.
var (
	userStore   = make(map[string]UserDetail)
	userStoreMu sync.RWMutex
)

// -----------------------------
// Handlers
// -----------------------------

// authHandler handles POST /auth
func authHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	// Decode JSON body into AuthRequestBody
	var reqBody AuthRequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "invalid JSON"})
		return
	}

	username := strings.TrimSpace(reqBody.Username)
	password := strings.TrimSpace(reqBody.Password)

	// Validate presence of username and password
	if len(username) == 0 || len(password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "username and password required"})
		return
	}

	// Placeholder authentication logic:
	// Only accept this exact pair. Replace with real logic as needed.
	if username == "testy.mctestface@example.com" && password == "Password1" {
		// Build a UserDetail object
		creationTime := time.Now().Unix()
		friendlyName := username
		if atIdx := strings.Index(username, "@"); atIdx > 0 {
			friendlyName = username[:atIdx]
		}
		encryptedPass := base64.StdEncoding.EncodeToString([]byte(password))

		userDetail := UserDetail{
			Username:          username,
			CreationDate:      creationTime,
			FriendlyName:      friendlyName,
			Email:             username,
			Source:            "basic",
			EncryptedPassword: encryptedPass,
			ExternalID:        username,
		}

		// Store (or overwrite) this UserDetail in our in-memory store
		userStoreMu.Lock()
		userStore[username] = userDetail
		userStoreMu.Unlock()

		// Respond with the external_user_identifier
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(AuthResponseBody{ExternalUserIdentifier: username})
		return
	}

	// Invalid credentials
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(ErrorResponse{Message: "invalid credentials"})
}

// getUserHandler handles GET /auth/users/{userid}
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	// Extract {userid} from the path.
	// The registered pattern is "/auth/users/", so anything after that is the userID.
	prefix := "/auth/users/"
	if !strings.HasPrefix(r.URL.Path, prefix) {
		http.NotFound(w, r)
		return
	}
	userID := strings.TrimPrefix(r.URL.Path, prefix)
	userID = strings.TrimSpace(userID)
	if userID == "" {
		http.NotFound(w, r)
		return
	}

	// Look up the user in our in-memory store
	userStoreMu.RLock()
	detail, exists := userStore[userID]
	userStoreMu.RUnlock()

	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "user not found"})
		return
	}

	// Return the UserDetail as JSON with HTTP 200
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(detail)
}

// -----------------------------
// Middleware
// -----------------------------

// loggingMiddleware prints out every incoming HTTP request (method, URI, remote addr, headers)
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log basic request line
		log.Printf("→ %s %s from %s", r.Method, r.RequestURI, r.RemoteAddr)
		// Log all request headers
		for name, values := range r.Header {
			for _, value := range values {
				log.Printf("    Header: %s: %s", name, value)
			}
		}
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

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
