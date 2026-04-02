## 📁 Environment Variables

This project uses environment variables loaded from a `.env` file via [`github.com/joho/godotenv`](https://github.com/joho/godotenv).

All configuration is centralized and validated on application startup.

---

### 📄 Setup

1. Copy the example file:

```bash
cp .env.example .env
```
2. Fill in all required variables in `.env`

> ⚠️ The .env file must be placed in the project root and must not be committed to git

### 🧾 Example .env

```bash
# Telegram MTProto API credentials
TG_API_ID=
TG_API_HASH=
PHONE_NUMBER=
PASSWORD=

# Telegram Bot API credentials
TG_BOT_TOKEN=
MODERATION_CHAT_ID=
TG_CHANNEL_ID=

# Sources for Telegram and RSS feeds
TG_SOURCES=
RSS_SOURCES=

# PostgreSQL database credentials
POSTGRES_DB=
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_PORT=
DATABASE_URL=
```

### 🔍 Notes

`TG_SOURCES` and `RSS_SOURCES` must be comma-separated:
```bash
TG_SOURCES=funny_memes,science_memes
RSS_SOURCES=https://www.reddit.com/r/TheMemeSub/.rss,https://www.reddit.com/r/RUSSIANMemeSub/.rss
```
