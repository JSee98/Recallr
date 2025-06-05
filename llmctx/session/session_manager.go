package session

import (
	"time"

	"github.com/JSee98/Recallr/models"
)

type SessionManager interface {
	// Create a new session for the given user with a TTL
	CreateSession(userID string, ttl time.Duration) (string, error)

	// Retrieve session metadata by ID
	GetSession(sessionID string) (*models.Session, error)

	// Gracefully end (delete) a session and clean up storage
	EndSession(sessionID string) error

	// Append a new message to the session's message list
	AddMessage(sessionID string, msg models.Message) error

	// Retrieve the last N messages from the session
	GetRecentMessages(sessionID string, limit int) ([]models.Message, error)

	// List all active session IDs for a user
	ListSessions(userID string) ([]string, error)
}
