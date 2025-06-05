package prompt

import (
	"context"
	"fmt"

	"github.com/JSee98/Recallr/memory"
	"github.com/JSee98/Recallr/models"
	"github.com/JSee98/Recallr/session"
)

type DefaultPromptBuilder struct {
	PromptManager  *PromptManager
	SessionManager session.SessionManager
	UserMemory     memory.UserMemory
}

func NewDefaultPromptBuilder(pm *PromptManager, sm session.SessionManager, um memory.UserMemory) *DefaultPromptBuilder {
	return &DefaultPromptBuilder{
		PromptManager:  pm,
		SessionManager: sm,
		UserMemory:     um,
	}
}

func (pb *DefaultPromptBuilder) BuildPrompt(ctx context.Context, sessionID, userID, currentInput string, messageLimit int) ([]models.Message, error) {
	messages := []models.Message{}

	// 1. Add system prompt
	messages = append(messages, models.Message{
		Role:    "system",
		Content: pb.PromptManager.SystemPrompt,
	})

	// 2. Inject long-term memory
	if facts, err := pb.UserMemory.ListFacts(userID); err == nil && len(facts) > 0 {
		content := formatFactsAsContext(facts)
		messages = append(messages, models.Message{
			Role:    "system",
			Content: content,
		})
	}

	// 3. Add recent session history
	recent, err := pb.SessionManager.GetRecentMessages(sessionID, messageLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to get session history: %w", err)
	}
	messages = append(messages, recent...)

	// 4. Add current user input
	messages = append(messages, models.Message{
		Role:    "user",
		Content: currentInput,
	})

	return messages, nil
}

func formatFactsAsContext(facts map[string]string) string {
	output := "The user has the following known facts:\n"
	for k, v := range facts {
		output += fmt.Sprintf("- %s: %s\n", k, v)
	}
	return output
}
