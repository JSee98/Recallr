package chat

import (
	"context"

	"github.com/JSee98/Recallr/models"
)

// Orchestrator defines the high-level interface to handle chat
// interactions including memory, prompt building, and streaming.
type Orchestrator interface {
	HandleUserInput(ctx context.Context, sessionID, userID, userInput string) (*StreamResult, error)
	HandleUserInputFull(ctx context.Context, sessionID, userID, userInput string) (*models.Message, error)
}
