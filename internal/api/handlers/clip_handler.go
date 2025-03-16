package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"strings"

	"twitchclipsearch/internal/database"
	"twitchclipsearch/internal/logger"
)

// ClipHandler handles HTTP requests related to clips
type ClipHandler struct {
	db *database.DB
}

// NewClipHandler creates a new instance of ClipHandler
func NewClipHandler(db *database.DB) *ClipHandler {
	return &ClipHandler{db: db}
}

// ClipResponse represents the JSON response for clip endpoints
type ClipResponse struct {
	ID           string    `json:"id"`
	StreamerName string    `json:"streamer_name"`
	Title        string    `json:"title"`
	URL          string    `json:"url"`
	CreatedAt    time.Time `json:"created_at"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// GetClips handles requests to retrieve clips
func (h *ClipHandler) GetClips(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get query parameters
	streamer := r.URL.Query().Get("streamer")
	limit := 50 // Default limit

	// Get clips from database
	clips, err := h.db.GetClips(streamer, limit)
	if err != nil {
		logger.Error("Failed to get clips", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to retrieve clips"})
		return
	}

	// Convert to response format
	response := make([]ClipResponse, len(clips))
	for i, clip := range clips {
		response[i] = ClipResponse{
			ID:           clip.ID,
			StreamerName: clip.StreamerName,
			Title:        clip.Title,
			URL:          clip.URL,
			CreatedAt:    clip.CreatedAt,
		}
	}

	json.NewEncoder(w).Encode(response)
}

// SearchClips handles clip search requests
func (h *ClipHandler) SearchClips(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get query parameters
	query := r.URL.Query().Get("q")
	if query == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Search query is required"})
		return
	}

	// Search clips in database
	clips, err := h.db.SearchClips(strings.ToLower(query))
	if err != nil {
		logger.Error("Failed to search clips", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to search clips"})
		return
	}

	// Convert to response format
	response := make([]ClipResponse, len(clips))
	for i, clip := range clips {
		response[i] = ClipResponse{
			ID:           clip.ID,
			StreamerName: clip.StreamerName,
			Title:        clip.Title,
			URL:          clip.URL,
			CreatedAt:    clip.CreatedAt,
		}
	}

	json.NewEncoder(w).Encode(response)
}
