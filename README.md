# Recallr Chat Context Example

This example demonstrates how to run an end-to-end chat interaction using the Recallr Golang library with:

- ✅ DragonflyDB as a Redis-compatible memory store
- 🧠 Long-term user memory
- 💬 Session context tracking
- ✍️ Customizable prompt injection
- 🧩 Pluggable LLM clients (example uses OpenAI-compatible client)

---

## 🔧 Prerequisites

1. **Install [Docker](https://docs.docker.com/get-docker/)**  
2. **Run DragonflyDB locally**

```bash
docker run -d --name dragonfly -p 6379:6379 docker.dragonflydb.io/dragonflydb/dragonfly
```

3. **Clone Recallr**

```bash
git clone https://github.com/Jsee98/Recallr
cd Recallr
```

4. **Run Example**

```bash
go run examples/main.go
```

---

## 📦 What It Does

- Starts a Dragonfly Redis instance.
- Initializes session manager, user memory, and prompt manager.
- Stores user facts into long-term memory.
- Builds a prompt using system prompt + user memory + chat history + current input.
- Sends the prompt to an LLM client (here: OpenAI-compatible streaming client).
- Streams the assistant's response line-by-line.
- Stores both user input and assistant output in session memory.

---

## 🧠 Message Roles Explained

| Role   | Description                                                                 |
|--------|-----------------------------------------------------------------------------|
| system | Sets the LLM's behavior, tone, and persona                                 |
| user   | Contains either a starting prompt or actual user input                     |
| assistant | LLM response, either full or streamed                                    |

---

## 🧙‍♂️ Prompt System

Prompt configuration is done via environment variables:

- `RECALLR_SYSTEM_PROMPT`: Instructions for the LLM’s behavior (required)
- `RECALLR_USER_PROMPT`: Optional user-level bootstrap message

These are injected during prompt building and can be reloaded at runtime:

```go
promptMgr.Reload()
```

---

## 🧠 Fact Summarizer

Recallr supports optional summarization of long-term user memory.

If your LLM client implements:

```go
type FactSummarizer interface {
    Summarizer(ctx context.Context, facts map[string]string) (string, error)
}
```

Then this summarizer is used to compress user facts before injecting into the prompt.

Otherwise, the default summarizer generates a simple readable block:

```
User Facts:
- location: Berlin
- language: Go
```

---

## 🔌 LLM Integration

This example includes an `OpenAICompatibleClient` which works with:

- OpenAI’s `chat/completions`
- DeepInfra, Fireworks, OpenRouter, Groq (any OpenAI-compatible proxy)

It supports both:
- `Chat(ctx, messages)` → full assistant reply
- `StreamChat(ctx, messages)` → streaming response via SSE

The orchestrator wraps these and exposes both:

```go
HandleUserInput(...)        // streaming (returns StreamResult)
HandleUserInputFull(...)    // full response
```

---

## 🧪 Testing & Extending

Want to add a custom LLM or summarizer? Just implement:

- `llm.Client` interface for chat/stream
- `prompt.FactSummarizer` for memory compression

Then wire them into the orchestrator like this:

```go
orchestrator := chat.NewDefaultOrchestrator(sessionMgr, promptBuilder, llmClient)
```

---

## 🧹 Cleanup

```bash
docker stop dragonfly && docker rm dragonfly
```

---

## 📜 License

MIT – see [LICENSE](./LICENSE)