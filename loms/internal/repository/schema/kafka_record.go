package schema

import (
	"time"
)

type KafkaRecordState int

const (
	PendingDelivery KafkaRecordState = iota
	Delivered
	MaxAttemptsReached
)

type KafkaRecord struct {
	ID           int              `db:"outbox_record_id"`
	Topic        string           `db:"topic"`
	PartitionKey string           `db:"key"`
	Message      string           `db:"message"`
	State        KafkaRecordState `db:"state"`
	CreatedOn    time.Time        `db:"created_on"`
}
