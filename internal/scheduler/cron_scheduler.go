package scheduler

import (
	"time"

	"github.com/robfig/cron/v3"
)

type CronScheduler struct {
	cron *cron.Cron
}

func NewCronScheduler() *CronScheduler {
	return &CronScheduler{
		cron: cron.New(cron.WithLocation(time.Local)),
	}
}

func (s *CronScheduler) RegisterJobs(
	ingestion func(),
	moderation func(),
	publish func(),
) error {

	// Ingestion: sunday 10:00
	if _, err := s.cron.AddFunc("0 10 * * 0", ingestion); err != nil {
		return err
	}

	// Send to moderation: daily 9:00 and 19:00
	if _, err := s.cron.AddFunc("0 9,19 * * *", moderation); err != nil {
		return err
	}

	// Publish: daily 10:00 and 20:00
	if _, err := s.cron.AddFunc("0 10,20 * * *", publish); err != nil {
		return err
	}

	return nil
}

func (s *CronScheduler) Start() {
	s.cron.Start()
}
