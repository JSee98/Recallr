# 🔁 Recallr

**Persistent LLM chat memory, session context, and global user state for Golang.**  
Drop-in memory infrastructure for building smart, context-aware LLM apps.

---

## 🚀 What is Recallr?

Recallr is a lightweight, pluggable Go library that enables:

- ✅ **Session memory** – track ongoing conversations
- ✅ **Chat recovery** – restore full context after a session ends
- ✅ **Global user memory** – build long-term user traits/preferences
- ✅ **DiceDB-backed persistence** – fully local, fast, and embeddable
- ✅ **Framework-free** – use with any LLM API or stack

No Python. No extra servers. Just Go.

---

## 🧠 Core Concepts

| Component         | Description                                                                 |
|------------------|-----------------------------------------------------------------------------|
| `SessionManager` | Manages in-session chat memory                                              |
| `ChatHistory`    | Stores and loads full chat histories per session                            |
| `UserMemory`     | Long-term memory store with user facts, traits, preferences                 |
| `PromptBuilder`  | Composes input prompts with memory, history, and current user message       |
| `LLMClient`      | Interface wrapper to plug in OpenAI, local models, or anything else         |

---

## 🛠️ Quick Start

```
import "github.com/your-org/recallr"

llm := recallr.New(recallr.Config{
    Storage: DiceDBStore,
    LLM:     OpenAIClient,
})

response := llm.Handle("user123", "session456", "What's the weather in Berlin?")
fmt.Println(response)
```
