package session

import (
	"context"
	"time"

	"github.com/Jsee98/Recallr/models"
)

type SessionStore interface {
	SaveSession(ctx context.Context, session *models.Session, ttl time.Duration) error
	GetSession(ctx context.Context, sessionID string) (*models.Session, error)
	DeleteSession(ctx context.Context, sessionID string) error

	AddMessage(ctx context.Context, sessionID string, msg *models.Message) error
	GetMessages(ctx context.Context, sessionID string, limit int) ([]*models.Message, error)

	AddToUserSessions(ctx context.Context, userID, sessionID string) error
	RemoveFromUserSessions(ctx context.Context, userID, sessionID string) error
	GetUserSessions(ctx context.Context, userID string) ([]string, error)
}
