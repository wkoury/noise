package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// fileHandler serves files with proper error handling
func fileHandler(filename string, contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if file exists
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		http.ServeFile(w, r, filename)
	}
}

// StartHTTPServer starts the HTTP server with graceful shutdown
func StartHTTPServer() {
	// Set up routes with error handling
	http.HandleFunc("/", fileHandler("internal/server/views/index.html", "text/html"))
	http.HandleFunc("/main.wasm", fileHandler("cmd/wasm/main.wasm", "application/wasm"))
	http.HandleFunc("/wasm_exec.js", fileHandler("internal/server/static/wasm_exec.js", "application/javascript"))

	// Add a health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if port[0] != ':' {
		port = ":" + port
	}

	log.Printf("Starting server on port %s", port)
	log.Printf("Server files:")
	log.Printf("  - Index: %s", filepath.Join("internal/server/views/index.html"))
	log.Printf("  - WASM: %s", filepath.Join("cmd/wasm/main.wasm"))
	log.Printf("  - JS: %s", filepath.Join("internal/server/static/wasm_exec.js"))

	// Create server
	server := &http.Server{
		Addr:         port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start: ", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give the server 30 seconds to shutdown gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
