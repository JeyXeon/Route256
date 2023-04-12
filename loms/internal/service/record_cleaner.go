package service

import (
	"context"
	"route256/libs/logger"
	"time"

	"go.uber.org/zap"
)

func (s *Service) RunRecordCleaner(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			err := s.RemoveExpiredMessages(ctx)
			if err != nil {
				logger.Error("records cleaning failed", zap.Error(err))
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (s *Service) RemoveExpiredMessages(ctx context.Context) error {
	maxRecordLifetime := 5 * time.Second
	expireTime := time.Now().UTC().Add(-1 * maxRecordLifetime)
	err := s.outboxKafkaRepository.RemoveRecordsBeforeDatetime(ctx, expireTime)
	if err != nil {
		return err
	}
	return nil
}
