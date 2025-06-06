package dragonfly

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Jsee98/Recallr/models"
	"github.com/Jsee98/Recallr/storage"
)

type DragonflySessionStore struct {
	store storage.Store
}

func NewDragonflySessionStore(store storage.Store) *DragonflySessionStore {
	return &DragonflySessionStore{store: store}
}

func (s *DragonflySessionStore) SaveSession(ctx context.Context, session *models.Session, ttl time.Duration) error {
	key := fmt.Sprintf("session:%s:meta", session.ID)

	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return s.store.Set(ctx, key, string(data), ttl)
}

func (s *DragonflySessionStore) GetSession(ctx context.Context, sessionID string) (*models.Session, error) {
	key := fmt.Sprintf("session:%s:meta", sessionID)

	data, err := s.store.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var session models.Session
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *DragonflySessionStore) DeleteSession(ctx context.Context, sessionID string) error {
	metaKey := fmt.Sprintf("session:%s:meta", sessionID)
	msgKey := fmt.Sprintf("session:%s:messages", sessionID)
	metaKeyErr := s.store.Delete(ctx, metaKey)
	msgKeyErr := s.store.Delete(ctx, msgKey)
	errToReturn := errors.Join(metaKeyErr, msgKeyErr)
	return errToReturn
}

func (s *DragonflySessionStore) AddMessage(ctx context.Context, sessionID string, msg *models.Message) error {
	key := fmt.Sprintf("session:%s:messages", sessionID)

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return s.store.LPush(ctx, key, string(data))
}

func (s *DragonflySessionStore) GetMessages(ctx context.Context, sessionID string, limit int) ([]*models.Message, error) {
	key := fmt.Sprintf("session:%s:messages", sessionID)

	// LRange returns newest messages last → limit from end
	raw, err := s.store.LRange(ctx, key, -int64(limit), -1)
	if err != nil {
		return nil, err
	}

	var msgs []*models.Message
	for _, item := range raw {
		var msg models.Message
		if err := json.Unmarshal([]byte(item), &msg); err == nil {
			msgs = append(msgs, &msg)
		}
	}
	return msgs, nil
}

func (s *DragonflySessionStore) AddToUserSessions(ctx context.Context, userID, sessionID string) error {
	key := fmt.Sprintf("user:%s:sessions", userID)
	return s.store.LPush(ctx, key, sessionID)
}

func (s *DragonflySessionStore) RemoveFromUserSessions(ctx context.Context, userID, sessionID string) error {
	// Simplified: you can rebuild this with a full list filter if needed
	// Dragonfly/Redis doesn’t have LREM in your wrapper yet
	return nil // to be implemented based on your use case
}

func (s *DragonflySessionStore) GetUserSessions(ctx context.Context, userID string) ([]string, error) {
	key := fmt.Sprintf("user:%s:sessions", userID)
	return s.store.LRange(ctx, key, 0, -1)
}
