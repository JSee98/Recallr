package chat

import (
	"context"
)

// Orchestrator defines the high-level interface to handle chat
// interactions including memory, prompt building, and streaming.
type Orchestrator interface {
	// HandleUserInput processes a user message and streams the assistant's response.
	// It returns two channels: one for streamed chunks and one for potential errors.
	HandleUserInput(ctx context.Context, sessionID, userID, userInput string) (<-chan string, <-chan error)
}
