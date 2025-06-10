package main

import "log"
import "net/http"

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