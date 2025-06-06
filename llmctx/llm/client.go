package llm

import (
	"context"
	"io"

	"github.com/JSee98/Recallr/models"
)

//go:generate mockgen -source=client.go -destination=../mocks/mock_llm_client.go -package=mocks
type Client interface {
	// Chat streams the LLM response based on prior messages.
	// Returns a reader you can consume line-by-line or chunk-by-chunk.
	StreamChat(ctx context.Context, messages []models.Message) (io.ReadCloser, error)

	// For simple use cases or testing: full message as one block.
	Chat(ctx context.Context, messages []models.Message) (*models.Message, error)

	// Name returns the LLM backend name (e.g., openai, llama-cpp, etc.)
	Name() string
}
