package services

import (
	"context"
	"fmt"
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
	fmt.Println("- processDueReminders", time.Now())
	return nil
}
