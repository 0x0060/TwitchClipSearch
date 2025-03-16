// Package service implements the core business logic for clip monitoring and processing
package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"twitchclipsearch/internal/config"
	"twitchclipsearch/internal/database"
	"twitchclipsearch/internal/discord"
	"twitchclipsearch/internal/logger"
	"twitchclipsearch/internal/metrics"

	"github.com/nicklaw5/helix/v2"
	"golang.org/x/time/rate"
)

// ClipService handles the core business logic for monitoring and processing Twitch clips
type ClipService struct {
	config     *config.Config
	db         *database.DB
	twitch     *helix.Client
	limiter    *rate.Limiter
	workerPool *WorkerPool
	shutdown   chan struct{}
	wg         sync.WaitGroup
}

// NewClipService creates a new instance of ClipService with the provided dependencies
func NewClipService(cfg *config.Config, db *database.DB) (*ClipService, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:     cfg.TwitchClientID,
		ClientSecret: cfg.TwitchClientSecret,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Twitch client: %w", err)
	}

	// Initialize rate limiter with Twitch API limits (30 requests per minute)
	limiter := rate.NewLimiter(rate.Every(2*time.Second), 1)

	// Create worker pool with configurable size
	pool := NewWorkerPool(5) // Adjust pool size based on needs

	return &ClipService{
		config:     cfg,
		db:         db,
		twitch:     client,
		limiter:    limiter,
		workerPool: pool,
		shutdown:   make(chan struct{}),
	}, nil
}

// Start begins the clip monitoring service
func (s *ClipService) Start(ctx context.Context) error {
	// Start the worker pool
	s.workerPool.Start()

	// Start monitoring for each streamer
	for streamerName := range s.config.Streamers {
		s.wg.Add(1)
		go s.monitorStreamer(ctx, streamerName)
	}

	return nil
}

// Stop gracefully shuts down the service
func (s *ClipService) Stop() error {
	// Signal shutdown
	close(s.shutdown)

	// Wait for all goroutines to finish
	s.wg.Wait()

	// Stop the worker pool
	s.workerPool.Stop()

	return nil
}

// monitorStreamer continuously monitors a streamer's clips
func (s *ClipService) monitorStreamer(ctx context.Context, streamerName string) {
	defer s.wg.Done()

	ticker := time.NewTicker(time.Duration(s.config.CheckIntervalSecs) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.shutdown:
			return
		case <-ticker.C:
			s.checkNewClips(ctx, streamerName)
		}
	}
}

// checkNewClips fetches and processes new clips for a streamer
func (s *ClipService) checkNewClips(ctx context.Context, streamerName string) {
	// Wait for rate limit
	err := s.limiter.Wait(ctx)
	if err != nil {
		logger.Error("Rate limit wait error", "error", err)
		return
	}

	// Get user ID for the streamer
	users, err := s.twitch.GetUsers(&helix.UsersParams{
		Logins: []string{streamerName},
	})
	if err != nil {
		logger.Error("Failed to get Twitch user", "error", err, "streamer", streamerName)
		metrics.RecordError("twitch_api_error")
		return
	}

	if len(users.Data.Users) == 0 {
		return
	}

	userID := users.Data.Users[0].ID

	// Get latest clip time from database
	latestTime, err := s.db.GetLatestClipTime(streamerName)
	if err != nil {
		logger.Error("Failed to get latest clip time", "error", err, "streamer", streamerName)
		metrics.RecordError("database_error")
		return
	}

	// Fetch clips after the latest time
	clips, err := s.twitch.GetClips(&helix.ClipsParams{
		BroadcasterID: userID,
		StartedAt:     latestTime,
	})
	if err != nil {
		logger.Error("Failed to fetch clips", "error", err, "streamer", streamerName)
		metrics.RecordError("twitch_api_error")
		return
	}

	// Process new clips using worker pool
	for _, clip := range clips.Data.Clips {
		clipData := clip // Create new variable to avoid closure issues
		s.workerPool.Submit(func() {
			s.processClip(ctx, streamerName, &clipData)
		})
	}
}

// processClip handles individual clip processing and storage
func (s *ClipService) processClip(ctx context.Context, streamerName string, clip *helix.Clip) {
	// Check if clip already exists
	exists, err := s.db.ClipExists(clip.ID)
	if err != nil {
		logger.Error("Failed to check clip existence", "error", err, "clip_id", clip.ID, "streamer", streamerName)
		metrics.RecordError("database_error")
		return
	}
	if exists {
		return
	}

	// Convert clip data
	createdAt, err := time.Parse(time.RFC3339, clip.CreatedAt)
	if err != nil {
		logger.Error("Failed to parse clip creation time", "error", err, "clip_id", clip.ID, "streamer", streamerName)
		metrics.RecordError("clip_processing_error")
		return
	}

	dbClip := &database.Clip{
		ID:           clip.ID,
		StreamerName: streamerName,
		Title:        clip.Title,
		URL:          clip.URL,
		CreatedAt:    createdAt,
		PostedAt:     time.Now(),
	}

	// Save to database
	if err := s.db.SaveClip(dbClip); err != nil {
		logger.Error("Failed to save clip", "error", err, "clip_id", clip.ID, "streamer", streamerName)
		metrics.RecordError("database_error")
		return
	}

	// Send notification (webhook)
	s.sendNotification(streamerName, dbClip)
}

// sendNotification sends a webhook notification for a new clip
func (s *ClipService) sendNotification(streamerName string, clip *database.Clip) {
	// Get webhook URL for the streamer
	webhookURL, ok := s.config.Streamers[streamerName]
	if !ok {
		return
	}

	// Create Discord client
	discordConfig := &discord.ClientConfig{
		WebhookURL: webhookURL,
		Username:   "TwitchClipBot",
		RateLimit:  5,
	}

	client := discord.NewClient(discordConfig)

	// Send notification
	if err := client.SendClipNotification(clip); err != nil {
		// Record metric for failed webhook
		metrics.RecordRateLimitHit("discord_webhook")
	}
}
