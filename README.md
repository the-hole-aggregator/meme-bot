# 🤖 Meme bot

Welcome to the Meme Aggregator Bot – a modular, scalable Go Telegram bot built using Clean Architecture, and a robust ingestion pipeline.

This bot collects memes from multiple sources – Telegram channels, RSS feeds – and processes them through a full-featured pipeline:
* **Fetch**: gather meme candidates from all configured sources
* **Validate**: ensure image format, size, and resolution meet requirements
* **Deduplicate**: generate image hashes to prevent reposts
* **Moderate**: send memes to a private moderation channel for review
* **Publish**: post approved memes to a public Telegram channel on a scheduled basis

---

## 📀 Architecture

The project architecture is fully structured around Clean Architecture with clear separation between layers: `entities`, `services`, `controllers`. The design encourages high testability, decoupling, and modularization.

📄 See detailed architecture overview: [docs/architecture.md](docs/architecture.md)

---

## 📚 Documentation

| Section                                                         | Description                              |
| --------------------------------------------------------------- | ---------------------------------------- |
| [architecture.md](docs/architecture.md)                         | Clean Architecture and layering overview |
| [git\_workflow\_guidelines.md](docs/git_workflow_guidelines.md) | Git flow and branching rules             |

---

## 🚀 Getting Started

1. **Configure git hooks**:

   ```bash
   make hooks-init
   ```

---

## 🧠 Philosophy

* Single responsibility at all levels
* Domain logic is independent of the realization details
