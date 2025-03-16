package discord

import (
	"fmt"
	"time"

	"twitchclipsearch/internal/database"
)

// Message represents a Discord webhook message
type Message struct {
	Username  string   `json:"username,omitempty"`
	Content   string   `json:"content,omitempty"`
	Embeds    []Embed  `json:"embeds,omitempty"`
}

// Embed represents a Discord message embed
type Embed struct {
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	URL         string    `json:"url"`
	Color       int       `json:"color"`
	Timestamp   string    `json:"timestamp,omitempty"`
	Fields      []Field   `json:"fields,omitempty"`
}

// Field represents a field in a Discord embed
type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// NewMessage creates a new Discord message from a clip
func NewMessage(clip *database.Clip) *Message {
	return &Message{
		Embeds: []Embed{
			{
				Title: clip.Title,
				Description: fmt.Sprintf("New clip from %s!", clip.StreamerName),
				URL: clip.URL,
				Color: 0x6441A4, // Twitch purple
				Timestamp: clip.CreatedAt.Format(time.RFC3339),
				Fields: []Field{
					{
						Name: "Streamer",
						Value: clip.StreamerName,
						Inline: true,
					},
					{
						Name: "Created At",
						Value: clip.CreatedAt.Format("Jan 02, 2006 15:04:05 MST"),
						Inline: true,
					},
				},
			},
		},
	}
}