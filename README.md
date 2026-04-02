# 🤖 Meme bot

Welcome to the Meme Aggregator Bot – Go Telegram bot built using Clean Architecture, and a robust ingestion pipeline.

This bot collects memes from multiple sources – Telegram channels, RSS feeds – and processes them through a full-featured pipeline:
* **Fetch**: gather meme candidates from all configured sources
* **Validate**: ensure image format, size, and resolution meet requirements
* **Deduplicate**: generate image hashes to prevent reposts
* **Moderate**: send memes to a private moderation channel for review
* **Publish**: post approved memes to a public Telegram channel on a scheduled basis

👉 [Follow link to subscribe](https://t.me/the_hole_memes)

---

## 📀 Architecture

The project architecture is based on Clean Architecture principles, with a clear separation of concerns across layers such as domain (entities), use_case (application business logic), adapters, delivery and scheduler.

#### Core Principles:
* Dependency Inversion: All dependencies point inward
* Isolation of business logic: The core logic does not depend on frameworks, databases, or external services
* Interface-driven design: business logic layer interacts with details via interfaces (ports)
* Testability: Business logic is easily testable using mocks

📄 See detailed architecture overview: [docs/architecture.md](docs/architecture.md)

---

## 📚 Documentation

| Section                                                         | Description                        |
| --------------------------------------------------------------- | ---------------------------------- |
| [architecture.md](docs/architecture.md)                         | Layers overview                    |
| [git\_workflow\_guidelines.md](docs/git_workflow_guidelines.md) | Git flow and branching rules       |
| [env.md](docs/env.md)                                           | Environment configs and .env usage |

---

## 🚀 Getting Started

1. **Configure git hooks**:

   ```bash
   make hooks-init
   ```
   
2. **Get dependencies**:

   ```bash
   go mod tidy
   ```
   
3. **Run the application:**

   ```bash
   make dev-up
   ```

---
