package kafka

import (
	"log"

	"github.com/Shopify/sarama"
)

type Consumer struct {
	ready chan bool
}

func NewConsumerGroup() Consumer {
	return Consumer{
		ready: make(chan bool),
	}
}

func (consumer *Consumer) Ready() <-chan bool {
	return consumer.ready
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
