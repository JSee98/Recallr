package memory

import (
	"context"
	"fmt"

	"github.com/Jsee98/Recallr/storage"
)

//go:generate mockgen -source=user_memory.go -destination=../mocks/mock_user_memory.go -package=mocks
type UserMemory interface {
	SetFact(userID, key, value string) error
	GetFact(userID, key string) (string, error)
	DeleteFact(userID, key string) error
	ListFacts(userID string) (map[string]string, error)
}

type DefaultUserMemory struct {
	store storage.Store
	ctx   context.Context
}

func NewUserMemory(store storage.Store) *DefaultUserMemory {
	return &DefaultUserMemory{
		store: store,
		ctx:   context.Background(),
	}
}

func (um *DefaultUserMemory) memoryKey(userID string) string {
	return fmt.Sprintf("user:%s:memory", userID)
}

func (um *DefaultUserMemory) SetFact(userID, key, value string) error {
	return um.store.HSet(um.ctx, um.memoryKey(userID), key, value)
}

func (um *DefaultUserMemory) GetFact(userID, key string) (string, error) {
	return um.store.HGet(um.ctx, um.memoryKey(userID), key)
}

func (um *DefaultUserMemory) DeleteFact(userID, key string) error {
	return um.store.HDel(um.ctx, um.memoryKey(userID), key)
}

func (um *DefaultUserMemory) ListFacts(userID string) (map[string]string, error) {
	return um.store.HGetAll(um.ctx, um.memoryKey(userID))
}
