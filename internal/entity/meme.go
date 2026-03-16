package entity

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
	ID        string
	Hash      string
	Source    Source
	SourceID  string
	Status    MemeStatus
	CreatedAt time.Time
}
