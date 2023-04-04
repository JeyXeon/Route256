package kafka

import (
	"github.com/Shopify/sarama"
)

type Listener interface {
	Subscribe() error
}

type HandleFunc func(message *sarama.ConsumerMessage)

type listener struct {
	consumer sarama.Consumer
	topic    string
	handler  HandleFunc
}

func NewListener(consumer sarama.Consumer, topic string, handler HandleFunc) Listener {
	return &listener{
		consumer: consumer,
		topic:    topic,
		handler:  handler,
	}
}

func (l *listener) Subscribe() error {
	partitionList, err := l.consumer.Partitions(l.topic)
	if err != nil {
		return err
	}

	for _, partition := range partitionList {
		initialOffset := sarama.OffsetOldest

		pc, err := l.consumer.ConsumePartition(l.topic, partition, initialOffset)
		if err != nil {
			return err
		}

		go func(pc sarama.PartitionConsumer) {
			for message := range pc.Messages() {
				l.handler(message)
			}
		}(pc)
	}

	return nil
}
