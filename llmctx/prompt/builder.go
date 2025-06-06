package prompt

import (
	"context"

	"github.com/JSee98/Recallr/models"
)

//go:generate mockgen -source=builder.go -destination=../mocks/mock_prompt_builder.go -package=mocks
type Builder interface {
	// BuildPrompt prepares the full LLM prompt using system prompt, memory, chat history, and current input.
	BuildPrompt(ctx context.Context, sessionID, userID, currentInput string, messageLimit int) ([]models.Message, error)
}
