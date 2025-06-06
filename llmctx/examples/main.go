package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/JSee98/Recallr/chat"
	"github.com/JSee98/Recallr/memory"
	"github.com/JSee98/Recallr/models"
	"github.com/JSee98/Recallr/prompt"
	"github.com/JSee98/Recallr/session"
	"github.com/JSee98/Recallr/storage/dragonfly"
)

type OpenAICompatibleClient struct {
	APIKey string
	APIURL string // e.g., https://api.openai.com/v1/chat/completions
	Model  string
}

func NewOpenAICompatibleClient(apiKey, apiURL, model string) *OpenAICompatibleClient {
	return &OpenAICompatibleClient{
		APIKey: apiKey,
		APIURL: apiURL,
		Model:  model,
	}
}

func (c *OpenAICompatibleClient) Name() string {
	return "openai"
}

func (c *OpenAICompatibleClient) Chat(ctx context.Context, messages []models.Message) (*models.Message, error) {
	reqBody := map[string]interface{}{
		"model":    c.Model,
		"messages": messages,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequestWithContext(ctx, "POST", c.APIURL, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Choices []struct {
			Message models.Message `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned")
	}
	return &result.Choices[0].Message, nil
}

func (c *OpenAICompatibleClient) StreamChat(ctx context.Context, messages []models.Message) (io.ReadCloser, error) {
	payload := map[string]interface{}{
		"model":    c.Model,
		"stream":   true,
		"messages": messages,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequestWithContext(ctx, "POST", c.APIURL, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		msg, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		return nil, fmt.Errorf("stream request failed: %s", msg)
	}

	return resp.Body, nil
}

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
	llmClient := NewOpenAICompatibleClient("YOUR_KEY", "https://api.deepinfra.com/v1/openai/chat/completions", "deepseek-ai/DeepSeek-V3")

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
	res, err := orchestrator.HandleUserInput(ctx, sessionID, userID, "What's the capital of France?")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Stream.Close()

	go func() {
		for err := range res.Errors {
			fmt.Println("stream error:", err)
		}
	}()

	go func() {
		for {
			buf := make([]byte, 256)
			n, err := res.Stream.Read(buf)
			if n > 0 {
				fmt.Print(string(buf[:n]))
			}
			if err != nil {
				break
			}
		}
	}()

	final := <-res.FinalMessage
	fmt.Println("full assistant message:", final.Content)
}
