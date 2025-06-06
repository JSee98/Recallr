package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/JSee98/Recallr/chat"
	"github.com/JSee98/Recallr/memory"
	"github.com/JSee98/Recallr/models"
	"github.com/JSee98/Recallr/prompt"
	"github.com/JSee98/Recallr/session"
	"github.com/JSee98/Recallr/storage/dragonfly"
)

type DummyLLMClient struct{}

func (d *DummyLLMClient) StreamChat(ctx context.Context, messages []models.Message) (io.ReadCloser, error) {
	pr, pw := io.Pipe()
	go func() {
		for _, msg := range messages {
			if msg.Role == "user" {
				pw.Write([]byte("Echo: " + msg.Content + "\n"))
			}
		}
		pw.Close()
	}()
	return pr, nil
}

func (d *DummyLLMClient) Chat(ctx context.Context, messages []models.Message) (*models.Message, error) {
	return &models.Message{
		ID:      "m1",
		Role:    "assistant",
		Content: "Echo response",
		Time:    time.Now(),
	}, nil
}

func (d *DummyLLMClient) Name() string { return "dummy" }

func main() {
	os.Setenv("RECALLR_SYSTEM_PROMPT", "You are a helpful assistant.")
	os.Setenv("RECALLR_USER_PROMPT", "Respond to the user input.")

	// 1. Dragonfly setup
	redisURL := "localhost:6379"
	drgCnfg := dragonfly.DragonflyConfig{
		Addr: redisURL,
	}
	redisStore := dragonfly.NewDragonflyStore(&drgCnfg)

	// 2. Dependencies
	sessionMgr := session.NewSessionManager(redisStore)
	userMem := memory.NewUserMemory(redisStore)
	promptMgr := prompt.NewPromptManager()
	builder := prompt.NewDefaultPromptBuilder(promptMgr, sessionMgr, userMem)
	llmClient := &DummyLLMClient{}

	// 3. Chat orchestrator
	orchestrator := chat.NewDefaultOrchestrator(sessionMgr, builder, llmClient)

	// 4. Simulate chat
	userID := "u123"
	sessionID, _ := sessionMgr.CreateSession(userID, 30*time.Minute)

	// Add memory
	_ = userMem.SetFact(userID, "location", "Berlin")
	_ = userMem.SetFact(userID, "language", "Go")

	// Handle input
	ctx := context.Background()
	output, errs := orchestrator.HandleUserInput(ctx, sessionID, userID, "What's my profile?")

	for {
		select {
		case line, ok := <-output:
			if !ok {
				return
			}
			fmt.Print(line)
		case err := <-errs:
			fmt.Println("Error:", err)
			return
		}
	}
}
