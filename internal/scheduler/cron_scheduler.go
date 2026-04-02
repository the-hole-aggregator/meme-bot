package scheduler

import (
	"time"

	"github.com/go-faster/errors"
	"github.com/robfig/cron/v3"
)

type CronScheduler struct {
	cron *cron.Cron
}

func NewCronScheduler() (*CronScheduler, error) {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err !=nil {
		return nil, errors.Wrap(err, "failed on cron initialization")
	}
	
	return &CronScheduler{
		cron: cron.New(cron.WithLocation(loc)),
	}, nil
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
