package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"twitchclipsearch/internal/config"
	"twitchclipsearch/internal/database"
	"twitchclipsearch/internal/logger"
	"twitchclipsearch/internal/metrics"
	"twitchclipsearch/internal/service"
)

func main() {
	cfg, err := config.LoadConfig()
	flag.Parse()

	// Initialize logger
	logger := logger.NewLogger()
	defer logger.Sync()

	// Load configuration
	cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize metrics
	metrics.InitMetrics()

	// Initialize database
	db, err := database.New(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create clip service
	clipService, err := service.NewClipService(cfg, db)
	if err != nil {
		log.Fatalf("Failed to create clip service: %v", err)
	}

	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start service
	if err := clipService.Start(ctx); err != nil {
		log.Fatalf("Failed to start service: %v", err)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Cleanup
	cancel()
	if err := clipService.Stop(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
}
