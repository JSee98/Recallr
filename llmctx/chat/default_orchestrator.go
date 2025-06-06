package chat

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"

	"github.com/JSee98/Recallr/llm"
	"github.com/JSee98/Recallr/models"
	"github.com/JSee98/Recallr/prompt"
	"github.com/JSee98/Recallr/session"
)

type DefaultOrchestrator struct {
	SessionManager session.SessionManager
	PromptBuilder  prompt.Builder
	LLMClient      llm.Client
}

func NewDefaultOrchestrator(sm session.SessionManager, pb prompt.Builder, lc llm.Client) *DefaultOrchestrator {
	return &DefaultOrchestrator{
		SessionManager: sm,
		PromptBuilder:  pb,
		LLMClient:      lc,
	}
}

func (o *DefaultOrchestrator) HandleUserInput(ctx context.Context, sessionID, userID, userInput string) (<-chan string, <-chan error) {
	out := make(chan string)
	errs := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errs)

		prompt, err := o.PromptBuilder.BuildPrompt(ctx, sessionID, userID, userInput, 10)
		if err != nil {
			errs <- fmt.Errorf("failed to build prompt: %w", err)
			return
		}

		userMsg := models.Message{
			ID:        uuid.New().String(),
			SessionID: sessionID,
			Role:      "user",
			Content:   userInput,
			Time:      time.Now(),
		}
		if err := o.SessionManager.AddMessage(sessionID, userMsg); err != nil {
			errs <- fmt.Errorf("failed to store user message: %w", err)
			return
		}

		stream, err := o.LLMClient.StreamChat(ctx, prompt)
		if err != nil {
			errs <- fmt.Errorf("stream error: %w", err)
			return
		}
		defer stream.Close()

		var fullResponse string
		reader := bufio.NewReader(stream)
		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				errs <- fmt.Errorf("read error: %w", err)
				return
			}
			fullResponse += line
			out <- line
		}

		llmMsg := models.Message{
			ID:        uuid.New().String(),
			SessionID: sessionID,
			Role:      "assistant",
			Content:   fullResponse,
			Time:      time.Now(),
		}
		if err := o.SessionManager.AddMessage(sessionID, llmMsg); err != nil {
			errs <- fmt.Errorf("failed to store LLM response: %w", err)
			return
		}
	}()

	return out, errs
}
