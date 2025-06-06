# 🐉 Testing DragonflyDB with Go Storage Layer

This guide helps you quickly spin up **DragonflyDB** in Docker and run unit tests against your Go storage interface implementation.

---

## 🚀 1. Run Dragonfly in Docker

```bash
# Pull the Dragonfly image
docker pull docker.dragonflydb.io/dragonflydb/dragonfly

# Run the container (exposes port 6379)
docker run -d \
  --name dragonfly \
  -p 6379:6379 \
  docker.dragonflydb.io/dragonflydb/dragonfly
