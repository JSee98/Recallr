package dragonfly

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/JSee98/Recallr/mocks"
	"github.com/JSee98/Recallr/models"
	"go.uber.org/mock/gomock"
)

func TestNewDragonflySessionStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	sessionStore := NewDragonflySessionStore(mockStore)

	if sessionStore == nil {
		t.Fatal("Expected non-nil session store")
	}
}

func TestSaveSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	sessionStore := NewDragonflySessionStore(mockStore)

	ctx := context.Background()
	session := &models.Session{
		ID:        "test-session-id",
		UserID:    "test-user-id",
		CreatedAt: time.Now().Truncate(time.Second),
		ExpiresAt: time.Now().Add(24 * time.Hour).Truncate(time.Second),
	}

	sessionData, _ := json.Marshal(session)

	// Expect the Set call with proper parameters
	mockStore.EXPECT().
		Set(ctx, "session:test-session-id:meta", string(sessionData), time.Hour).
		Return(nil)

	err := sessionStore.SaveSession(ctx, session, 1*time.Hour)
	if err != nil {
		t.Fatalf("SaveSession failed: %v", err)
	}
}

func TestSaveSessionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	sessionStore := NewDragonflySessionStore(mockStore)

	ctx := context.Background()
	session := &models.Session{
		ID:        "test-session-id",
		UserID:    "test-user-id",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	expectedErr := errors.New("set error")

	// Expect the Set call to return an error
	mockStore.EXPECT().
		Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(expectedErr)

	err := sessionStore.SaveSession(ctx, session, 1*time.Hour)
	if err == nil {
		t.Fatal("Expected error from SaveSession, got nil")
	}
	if err != expectedErr {
		t.Fatalf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestGetSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	sessionStore := NewDragonflySessionStore(mockStore)

	ctx := context.Background()
	originalSession := &models.Session{
		ID:        "test-session-id",
		UserID:    "test-user-id",
		CreatedAt: time.Now().Truncate(time.Second),
		ExpiresAt: time.Now().Add(24 * time.Hour).Truncate(time.Second),
	}

	sessionData, _ := json.Marshal(originalSession)

	// Expect the Get call with the proper key
	mockStore.EXPECT().
		Get(ctx, "session:test-session-id:meta").
		Return(string(sessionData), nil)

	session, err := sessionStore.GetSession(ctx, "test-session-id")
	if err != nil {
		t.Fatalf("GetSession failed: %v", err)
	}

	if session == nil {
		t.Fatal("Expected non-nil session")
	}

	if session.ID != originalSession.ID ||
		session.UserID != originalSession.UserID ||
		!session.CreatedAt.Equal(originalSession.CreatedAt) ||
		!session.ExpiresAt.Equal(originalSession.ExpiresAt) {
		t.Errorf("Retrieved session doesn't match original: got %+v, want %+v", session, originalSession)
	}
}

func TestGetSessionNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	sessionStore := NewDragonflySessionStore(mockStore)

	ctx := context.Background()
	expectedErr := errors.New("key not found")

	mockStore.EXPECT().
		Get(ctx, "session:non-existent-id:meta").
		Return("", expectedErr)

	_, err := sessionStore.GetSession(ctx, "non-existent-id")
	if err == nil {
		t.Fatal("Expected error when getting non-existent session, got nil")
	}
}

func TestGetSessionInvalidData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	sessionStore := NewDragonflySessionStore(mockStore)

	ctx := context.Background()

	// Return invalid JSON
	mockStore.EXPECT().
		Get(ctx, "session:invalid-json:meta").
		Return("not-valid-json", nil)

	_, err := sessionStore.GetSession(ctx, "invalid-json")
	if err == nil {
		t.Fatal("Expected error when parsing invalid session data, got nil")
	}
}

func TestDeleteSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	sessionStore := NewDragonflySessionStore(mockStore)

	ctx := context.Background()

	// Expect Delete calls for both keys
	mockStore.EXPECT().
		Delete(ctx, "session:test-id:meta").
		Return(nil)
	mockStore.EXPECT().
		Delete(ctx, "session:test-id:messages").
		Return(nil)

	err := sessionStore.DeleteSession(ctx, "test-id")
	if err != nil {
		t.Fatalf("DeleteSession failed: %v", err)
	}
}

func TestDeleteSessionPartialError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	sessionStore := NewDragonflySessionStore(mockStore)

	ctx := context.Background()
	expectedErr := errors.New("delete error")

	// Expect Delete calls with one returning an error
	mockStore.EXPECT().
		Delete(ctx, "session:test-id:meta").
		Return(expectedErr)
	mockStore.EXPECT().
		Delete(ctx, "session:test-id:messages").
		Return(nil)

	err := sessionStore.DeleteSession(ctx, "test-id")
	if err == nil {
		t.Fatal("Expected error from DeleteSession, got nil")
	}
}

func TestAddMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	sessionStore := NewDragonflySessionStore(mockStore)

	ctx := context.Background()
	msg := &models.Message{
		Role:    "user",
		Content: "Hello world",
		Time:    time.Now().Truncate(time.Second),
	}

	msgData, _ := json.Marshal(msg)

	// Expect LPush call with the serialized message
	mockStore.EXPECT().
		LPush(ctx, "session:test-session:messages", string(msgData)).
		Return(nil)

	err := sessionStore.AddMessage(ctx, "test-session", msg)
	if err != nil {
		t.Fatalf("AddMessage failed: %v", err)
	}
}

func TestGetMessages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	sessionStore := NewDragonflySessionStore(mockStore)

	ctx := context.Background()

	// Create test messages
	msgs := []*models.Message{
		{Role: "user", Content: "Hello", Time: time.Now().Add(-3 * time.Minute)},
		{Role: "assistant", Content: "Hi there", Time: time.Now().Add(-2 * time.Minute)},
		{Role: "user", Content: "How are you?", Time: time.Now().Add(-1 * time.Minute)},
	}

	// Serialize messages to mock return value
	serializedMsgs := make([]string, 0, len(msgs))
	for _, msg := range msgs {
		data, _ := json.Marshal(msg)
		serializedMsgs = append(serializedMsgs, string(data))
	}

	// Return last 2 messages
	mockStore.EXPECT().
		LRange(ctx, "session:test-session:messages", int64(-2), int64(-1)).
		Return(serializedMsgs[1:], nil)

	retrieved, err := sessionStore.GetMessages(ctx, "test-session", 2)
	if err != nil {
		t.Fatalf("GetMessages failed: %v", err)
	}

	if len(retrieved) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(retrieved))
	}

	if retrieved[0].Content != msgs[1].Content || retrieved[1].Content != msgs[2].Content {
		t.Errorf("Retrieved messages with wrong content")
	}
}

func TestAddToUserSessions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	sessionStore := NewDragonflySessionStore(mockStore)

	ctx := context.Background()

	mockStore.EXPECT().
		LPush(ctx, "user:user123:sessions", "session456").
		Return(nil)

	err := sessionStore.AddToUserSessions(ctx, "user123", "session456")
	if err != nil {
		t.Fatalf("AddToUserSessions failed: %v", err)
	}
}

func TestGetUserSessions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	sessionStore := NewDragonflySessionStore(mockStore)

	ctx := context.Background()
	expectedSessions := []string{"session1", "session2", "session3"}

	mockStore.EXPECT().
		LRange(ctx, "user:user123:sessions", int64(0), int64(-1)).
		Return(expectedSessions, nil)

	sessions, err := sessionStore.GetUserSessions(ctx, "user123")
	if err != nil {
		t.Fatalf("GetUserSessions failed: %v", err)
	}

	if len(sessions) != 3 {
		t.Fatalf("Expected 3 sessions, got %d", len(sessions))
	}

	// Check all sessions are present
	for i, sid := range sessions {
		if sid != expectedSessions[i] {
			t.Errorf("Session at index %d: expected '%s', got '%s'", i, expectedSessions[i], sid)
		}
	}
}

func TestGetUserSessionsEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	sessionStore := NewDragonflySessionStore(mockStore)

	ctx := context.Background()

	mockStore.EXPECT().
		LRange(ctx, "user:nonexistent:sessions", int64(0), int64(-1)).
		Return([]string{}, nil)

	sessions, err := sessionStore.GetUserSessions(ctx, "nonexistent")
	if err != nil {
		t.Fatalf("GetUserSessions failed: %v", err)
	}

	if len(sessions) != 0 {
		t.Errorf("Expected empty list for nonexistent user, got %d sessions", len(sessions))
	}
}
