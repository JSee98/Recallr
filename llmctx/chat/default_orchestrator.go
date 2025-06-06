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
	"github.com/JSee98/Recallr/types"
)

type StreamResult struct {
	Stream       io.ReadCloser
	Errors       <-chan error
	FinalMessage <-chan *models.Message
}

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

func (o *DefaultOrchestrator) HandleUserInput(ctx context.Context, sessionID, userID, userInput string) (*StreamResult, error) {
	// Select summarizer
	var summarizerFn types.SummarizerFunction = prompt.DefaultSummarizer
	if summarizer, ok := o.LLMClient.(prompt.FactSummarizer); ok {
		summarizerFn = summarizer.Summarizer
	}

	// Build prompt
	promptToSend, err := o.PromptBuilder.BuildPrompt(ctx, sessionID, userID, userInput, 10, summarizerFn)
	if err != nil {
		return nil, fmt.Errorf("failed to build prompt: %w", err)
	}

	// Add user message to session
	userMsg := models.Message{
		ID:        uuid.New().String(),
		SessionID: sessionID,
		Role:      "user",
		Content:   userInput,
		Time:      time.Now(),
	}
	if err := o.SessionManager.AddMessage(sessionID, userMsg); err != nil {
		return nil, fmt.Errorf("failed to store user message: %w", err)
	}

	// Stream from LLM client
	llmStream, err := o.LLMClient.StreamChat(ctx, promptToSend)
	if err != nil {
		return nil, fmt.Errorf("LLM stream error: %w", err)
	}

	// Setup pipe to stream back to caller
	pr, pw := io.Pipe()
	errs := make(chan error, 1)
	finalMsg := make(chan *models.Message, 1)

	go func() {
		defer llmStream.Close()
		defer pw.Close()
		defer close(errs)
		defer close(finalMsg)

		reader := bufio.NewReader(llmStream)
		var fullResponse string

		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				errs <- fmt.Errorf("stream read error: %w", err)
				return
			}
			fullResponse += line
			if _, err := pw.Write([]byte(line)); err != nil {
				errs <- fmt.Errorf("pipe write error: %w", err)
				return
			}
		}

		// Store final assistant message
		llmMsg := &models.Message{
			ID:        uuid.New().String(),
			SessionID: sessionID,
			Role:      "assistant",
			Content:   fullResponse,
			Time:      time.Now(),
		}
		if err := o.SessionManager.AddMessage(sessionID, *llmMsg); err != nil {
			errs <- fmt.Errorf("failed to store LLM response: %w", err)
			return
		}
		finalMsg <- llmMsg
	}()

	return &StreamResult{
		Stream:       pr,
		Errors:       errs,
		FinalMessage: finalMsg,
	}, nil
}

func (o *DefaultOrchestrator) HandleUserInputFull(ctx context.Context, sessionID, userID, userInput string) (*models.Message, error) {
	// Choose summarizer
	var summarizerFn types.SummarizerFunction = prompt.DefaultSummarizer
	if summarizer, ok := o.LLMClient.(prompt.FactSummarizer); ok {
		summarizerFn = summarizer.Summarizer
	}

	// Build prompt
	promptToSend, err := o.PromptBuilder.BuildPrompt(ctx, sessionID, userID, userInput, 10, summarizerFn)
	if err != nil {
		return nil, fmt.Errorf("failed to build prompt: %w", err)
	}

	// Store user message
	userMsg := models.Message{
		ID:        uuid.New().String(),
		SessionID: sessionID,
		Role:      "user",
		Content:   userInput,
		Time:      time.Now(),
	}
	if err := o.SessionManager.AddMessage(sessionID, userMsg); err != nil {
		return nil, fmt.Errorf("failed to store user message: %w", err)
	}

	// Get full response
	llmResp, err := o.LLMClient.Chat(ctx, promptToSend)
	if err != nil {
		return nil, fmt.Errorf("LLM call failed: %w", err)
	}

	llmResp.ID = uuid.New().String()
	llmResp.SessionID = sessionID
	llmResp.Role = "assistant"
	llmResp.Time = time.Now()

	// Store LLM response
	if err := o.SessionManager.AddMessage(sessionID, *llmResp); err != nil {
		return nil, fmt.Errorf("failed to store LLM response: %w", err)
	}

	return llmResp, nil
}
