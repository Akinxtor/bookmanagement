package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Akinxtor/bookmanagement/pkg/routes"
	"github.com/gorilla/mux"
)

// Middleware to handle CORS
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Create router and register routes
	r := mux.NewRouter()
	routes.RegisterBookStoreRoutes(r)

	// Health check
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "✅ Bookstore API is running")
	})

	// ✅ Wrap router with CORS middleware
	handlerWithCORS := enableCORS(r)

	// ✅ Pass the wrapped router to ListenAndServe
	fmt.Println("🚀 Bookstore API running at: http://localhost:9010/")
	log.Fatal(http.ListenAndServe("localhost:9010", handlerWithCORS))
}
