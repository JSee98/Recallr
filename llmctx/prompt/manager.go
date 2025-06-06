package prompt

import (
	"fmt"
	"os"
)

//go:generate mockgen -source=manager.go -destination=../mocks/mock_prompt_manager.go -package=mocks
type PromptManagerInterface interface{
	// Reload prompts from environment variables.
	Reload() error
}

// PromptManager holds all runtime prompts used by the LLM service.
type PromptManager struct {
	SystemPrompt string
	UserPrompt   string
	// Add more prompt types as needed (e.g., SummaryPrompt, FactPrompt)
}

// RequiredPromptEnvVars maps env variable names to human-friendly labels for error reporting.
var RequiredPromptEnvVars = map[string]*string{
	"RECALLR_SYSTEM_PROMPT": nil,
	"RECALLR_USER_PROMPT":   nil,
}

// NewPromptManager initializes the prompt manager from environment variables.
// It panics on missing required env vars to ensure safety at startup.
func NewPromptManager() *PromptManager {
	loadRequiredPrompts()

	return &PromptManager{
		SystemPrompt: *RequiredPromptEnvVars["RECALLR_SYSTEM_PROMPT"],
		UserPrompt:   *RequiredPromptEnvVars["RECALLR_USER_PROMPT"],
	}
}

// loadRequiredPrompts validates and loads each required env variable.
func loadRequiredPrompts() {
	for key := range RequiredPromptEnvVars {
		val := os.Getenv(key)
		if val == "" {
			panic(fmt.Sprintf("Missing required environment variable: %s", key))
		}
		RequiredPromptEnvVars[key] = &val
	}
}

// Reload can be called to hot-reload prompts from environment variables during runtime.
// Useful for long-running services with config changes.
func (pm *PromptManager) Reload() error {
	for key := range RequiredPromptEnvVars {
		val := os.Getenv(key)
		if val == "" {
			return fmt.Errorf("cannot reload: missing env var %s", key)
		}
		switch key {
		case "RECALLR_SYSTEM_PROMPT":
			pm.SystemPrompt = val
		case "RECALLR_USER_PROMPT":
			pm.UserPrompt = val
		}
	}
	return nil
}
