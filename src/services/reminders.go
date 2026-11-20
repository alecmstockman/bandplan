package services

import (
	"context"
	"log/slog"
	"time"
)

func (s *Service) RunReminderWorker(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.processDueReminders(ctx); err != nil {
				slog.Error("process due reminders", "error", err)
			}
		}
	}
}

func (s *Service) processDueReminders(ctx context.Context) error {
	return nil
}
