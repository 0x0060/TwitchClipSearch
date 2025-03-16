package database

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Clip represents a Twitch clip in the database
type Clip struct {
	ID          string
	StreamerName string
	Title        string
	URL          string
	CreatedAt    time.Time
	PostedAt     time.Time
}

// DB handles database operations
type DB struct {
	db *sql.DB
}

// New creates a new database connection and initializes the schema
func New(path string) (*DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err := initSchema(db); err != nil {
		db.Close()
		return nil, err
	}

	return &DB{db: db}, nil
}

// initSchema creates the necessary database tables if they don't exist
func initSchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS clips (
			id TEXT PRIMARY KEY,
			streamer_name TEXT NOT NULL,
			title TEXT NOT NULL,
			url TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			posted_at DATETIME NOT NULL
		)
	`)
	return err
}

// Close closes the database connection
func (d *DB) Close() error {
	return d.db.Close()
}

// SaveClip saves a new clip to the database
func (d *DB) SaveClip(clip *Clip) error {
	_, err := d.db.Exec(
		"INSERT INTO clips (id, streamer_name, title, url, created_at, posted_at) VALUES (?, ?, ?, ?, ?, ?)",
		clip.ID,
		clip.StreamerName,
		clip.Title,
		clip.URL,
		clip.CreatedAt,
		clip.PostedAt,
	)
	return err
}

// GetLatestClipTime returns the creation time of the most recent clip for a streamer
func (d *DB) GetLatestClipTime(streamerName string) (time.Time, error) {
	var createdAt time.Time
	err := d.db.QueryRow(
		"SELECT created_at FROM clips WHERE streamer_name = ? ORDER BY created_at DESC LIMIT 1",
		streamerName,
	).Scan(&createdAt)

	if err == sql.ErrNoRows {
		return time.Time{}, nil
	}
	return createdAt, err
}

// ClipExists checks if a clip with the given ID already exists
func (d *DB) ClipExists(clipID string) (bool, error) {
	var exists bool
	err := d.db.QueryRow("SELECT EXISTS(SELECT 1 FROM clips WHERE id = ?)", clipID).Scan(&exists)
	return exists, err
}