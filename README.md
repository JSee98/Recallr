# Recallr Chat Context Example

This example demonstrates how to run an end-to-end chat interaction using the Recallr Golang library with:

- DragonflyDB as a Redis-compatible memory store
- Long-term user memory
- Session context tracking
- Customizable prompt injection
- Pluggable LLM clients (example uses a dummy echo client)

---

## 🔧 Prerequisites

1. **Install [Docker](https://docs.docker.com/get-docker/)**
2. **Run DragonflyDB locally**

```bash
docker run -d --name dragonfly -p 6379:6379 docker.dragonflydb.io/dragonflydb/dragonfly
```

3. **Clone Recallr**

```bash
git clone https://github.com/JSee98/Recallr
cd Recallr
```

4. **Run Example**

```bash
go run examples/main.go
```

---

## 📦 What It Does

- Starts a Dragonfly Redis instance.
- Sets up session manager, user memory, and prompt manager.
- Stores user facts into long-term memory.
- Builds a prompt using the memory + last chat.
- Sends the prompt to an LLM client (dummy echo used here).
- Streams the assistant's response back line-by-line.
- Stores the exchange in session memory.

---

## 💡 Prompt System

Prompt variables are managed using environment variables:

- `RECALLR_SYSTEM_PROMPT`: System-level instructions.
- `RECALLR_USER_PROMPT`: Optional user prompt template.

You can hot-reload these by calling `.Reload()` on `PromptManager`.

---

## 🧠 Fact Summarizer

If the LLM client implements the following interface:

```go
type FactSummarizer interface {
    Summarizer(ctx context.Context, facts map[string]string) (string, error)
}
```

The orchestrator will use this to summarize user memory context into a compact string before sending the chat prompt.

Otherwise, a default summarizer will generate:

```
User Facts:
- key1: value1
- key2: value2
```

---

## 🧪 Testing & Extending

Replace the `DummyLLMClient` with your own `llm.Client` implementation (e.g., OpenAI, Claude, etc.) to integrate real LLMs.

---

## 🧹 Cleanup

```bash
docker stop dragonfly && docker rm dragonfly
```

---

## 📜 License

MIT – see [LICENSE](./LICENSE)
