package model

import (
	"fmt"
	desc "route256/loms/pkg/notifications"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
)

type RecordState int

const (
	PendingDelivery RecordState = iota
	Delivered
	MaxAttemptsReached
)

type KafkaRecord struct {
	ID           int
	Topic        string
	PartitionKey string
	Message      string
	State        RecordState
	CreatedOn    time.Time
}

func NewOrderStatusChangeKafkaRecord(orderId int64, orderStatus OrderStatus) (*KafkaRecord, error) {
	orderStatusChange := desc.OrderStatusChange{
		Id:     orderId,
		Status: string(orderStatus),
	}

	bytes, err := protojson.Marshal(&orderStatusChange)
	if err != nil {
		return nil, err
	}

	return &KafkaRecord{
		Topic:        "order_state_change",
		PartitionKey: fmt.Sprint(orderId),
		Message:      string(bytes),
		State:        PendingDelivery,
		CreatedOn:    time.Now(),
	}, nil
}
