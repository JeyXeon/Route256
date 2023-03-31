package producer

import (
	"fmt"
	"log"
	"route256/loms/internal/model"
	desc "route256/loms/pkg/notifications"
	"time"

	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/encoding/protojson"
)

type OrderStatusChangeProducer interface {
	SendOrderStatusChange(orderId int64, status model.OrderStatus) error
}

type orderStatusChangeProducer struct {
	producer sarama.AsyncProducer
	topic    string
}

func NewOrderStatusProducer(producer sarama.AsyncProducer, topic string) OrderStatusChangeProducer {
	orderStatusChangeProducer := &orderStatusChangeProducer{
		producer: producer,
		topic:    topic,
	}

	go func() {
		for e := range producer.Errors() {
			bytes, _ := e.Msg.Key.Encode()
			log.Printf("failed to process order: %s, error: %s", string(bytes), e.Error())
		}
	}()

	go func() {
		for m := range producer.Successes() {
			bytes, _ := m.Key.Encode()
			log.Printf("order id: %s, partition: %d, offset: %d", string(bytes), m.Partition, m.Offset)
		}
	}()

	return orderStatusChangeProducer

}

func (p *orderStatusChangeProducer) SendOrderStatusChange(orderId int64, status model.OrderStatus) error {
	orderStatusChange := desc.OrderStatusChange{
		Id:     orderId,
		Status: string(status),
	}

	bytes, err := protojson.Marshal(&orderStatusChange)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic:     p.topic,
		Partition: -1,
		Value:     sarama.ByteEncoder(bytes),
		Key:       sarama.StringEncoder(fmt.Sprint(orderId)),
		Timestamp: time.Now(),
	}

	p.producer.Input() <- msg

	return nil
}
