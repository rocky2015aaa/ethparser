package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rocky2015aaa/ethparser/handler"
	"github.com/rocky2015aaa/ethparser/service"
)

const (
	envPollInterval = "POLL_INTERVAL"
	defaultInterval = 5
)

//	@title			Ethereum Parser
//	@version		1.0
//	@description	Ethereum Parser

// @host		localhost:8080
// @BasePath	/
func main() {
	// Load configuration
	rpcURL := os.Getenv(service.EnvRpcUrl)
	if rpcURL == "" {
		err := os.Setenv(service.EnvRpcUrl, service.DefaultRpcUrl)
		if err != nil {
			log.Fatal("Error setting environment variable:", err)
		}
	}

	pollInterval, err := time.ParseDuration(os.Getenv(envPollInterval))
	if err != nil || pollInterval == 0 {
		pollInterval = defaultInterval * time.Second
	}

	// Initialize store and services
	blockService := service.NewBlockService()
	transactionHandler := handler.NewTransactionHandler(blockService)

	// Start loading blocks in the background
	go blockService.LoadBlocks(pollInterval)

	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	// Serve Swagger JSON (OpenAPI spec)
	http.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.json")
	})
	// Serve Swagger YAML (OpenAPI spec)
	http.HandleFunc("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.yaml")
	})
	// Serve Swagger UI HTML
	http.HandleFunc("/openapi.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/openapi.html")
	})
	http.HandleFunc("/subscribe", transactionHandler.HandleSubscribe)
	http.HandleFunc("/transactions", transactionHandler.HandleGetTransactions)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Server running on port 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	// Context with timeout to allow ongoing tasks to finish
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shut down the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
