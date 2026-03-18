package domain

import "time"

type MemeStatus string

const (
	Pending  MemeStatus = "pending"
	Approved MemeStatus = "approved"
	Posted   MemeStatus = "posted"
)

type Source string

const (
	Telegram Source = "telegram"
	RSS      Source = "rss"
)

type Meme struct {
	ID        int
	PHash     string
	Status    MemeStatus
	Source    Source
	SourceID  string
	CreatedAt time.Time
}
