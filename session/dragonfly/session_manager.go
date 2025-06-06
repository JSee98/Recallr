package dragonfly

import (
	"context"
	"fmt"
	"time"

	"github.com/JSee98/Recallr/models"
)

type SessionManager struct {
	storage *DragonflySessionStore
	ctx     context.Context
}

func NewSessionManager(storage *DragonflySessionStore) *SessionManager {
	return &SessionManager{
		storage: storage,
		ctx:     context.Background(),
	}
}

func (sm *SessionManager) CreateSession(userID string, ttl time.Duration) (string, error) {
	sessionID := fmt.Sprintf("%s:%d", userID, time.Now().UnixNano())

	session := &models.Session{
		ID:        sessionID,
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(ttl),
	}

	if err := sm.storage.SaveSession(sm.ctx, session, ttl); err != nil {
		return "", err
	}

	if err := sm.storage.AddToUserSessions(sm.ctx, userID, sessionID); err != nil {
		return "", err
	}

	return sessionID, nil
}

func (sm *SessionManager) GetSession(sessionID string) (*models.Session, error) {
	return sm.storage.GetSession(sm.ctx, sessionID)
}

func (sm *SessionManager) EndSession(sessionID string) error {
	session, err := sm.storage.GetSession(sm.ctx, sessionID)
	if err != nil {
		return err
	}

	if err := sm.storage.DeleteSession(sm.ctx, sessionID); err != nil {
		return err
	}

	return sm.storage.RemoveFromUserSessions(sm.ctx, session.UserID, sessionID)
}

func (sm *SessionManager) AddMessage(sessionID string, msg models.Message) error {
	return sm.storage.AddMessage(sm.ctx, sessionID, &msg)
}

func (sm *SessionManager) GetRecentMessages(sessionID string, limit int) ([]models.Message, error) {
	rawMsgs, err := sm.storage.GetMessages(sm.ctx, sessionID, limit)
	if err != nil {
		return nil, err
	}

	messages := make([]models.Message, len(rawMsgs))
	for i, m := range rawMsgs {
		messages[i] = *m
	}
	return messages, nil
}

func (sm *SessionManager) ListSessions(userID string) ([]string, error) {
	return sm.storage.GetUserSessions(sm.ctx, userID)
}
