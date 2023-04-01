package service

import (
	"context"
	"log"
	"route256/loms/internal/model"
	"time"

	"github.com/Shopify/sarama"
)

func (s *Service) RunRecordProcessor(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			err := s.transactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
				err := s.ProcessRecords(ctx)
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				log.Println(err.Error())
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (s *Service) ProcessRecords(ctx context.Context) error {
	records, err := s.outboxKafkaRepository.GetUnprocessedRecords(ctx)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return nil
	}

	return s.publishMessages(ctx, records)
}

func (s *Service) publishMessages(ctx context.Context, records []*model.KafkaRecord) error {
	for _, rec := range records {

		attempts := 5
		var publishErr error
		for attempts > 0 {
			err := s.publishMessage(rec)
			if err != nil {
				publishErr = err
			} else {
				break
			}

			attempts--
		}

		if publishErr != nil {
			rec.State = model.MaxAttemptsReached

			err := s.outboxKafkaRepository.UpdateRecordByID(ctx, rec)
			if err != nil {
				return err
			}

			return publishErr
		}

		rec.State = model.Delivered
		err := s.outboxKafkaRepository.UpdateRecordByID(ctx, rec)

		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) publishMessage(rec *model.KafkaRecord) error {
	msg := &sarama.ProducerMessage{
		Topic:     rec.Topic,
		Partition: -1,
		Value:     sarama.StringEncoder(rec.Message),
		Key:       sarama.StringEncoder(rec.PartitionKey),
		Timestamp: rec.CreatedOn,
	}

	_, _, err := s.kafkaSyncProducer.SendMessage(msg)

	return err
}
