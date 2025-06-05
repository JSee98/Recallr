package chat

import (
	"time"

	"github.com/JSee98/Recallr/models"
)

type ChatHistoryStore interface {
	SaveMessage(msg models.Message) error
	GetMessages(sessionID string, limit int, reverse bool) ([]models.Message, error)
	GetRecentMessages(sessionID string, since time.Time) ([]models.Message, error)
	DeleteMessages(sessionID string) error
}
