package main

import (
	"context"
	"log"
	"route256/libs/kafka"
	"route256/libs/logger"
	"route256/notifications/internal/config"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

func main() {
	logger.Init(false)

	err := config.Init()
	if err != nil {
		logger.Fatal("config init", zap.Error(err))
	}

	brokers := config.ConfigData.Kafka.Brokers
	c, err := kafka.NewConsumer(brokers)
	if err != nil {
		logger.Fatal("kafka consumer", zap.Error(err))
	}

	orderStateChangeTopicName := config.ConfigData.Kafka.Topics.OrderStateChange.Name
	orderStateChangeListener := kafka.NewListener(c, orderStateChangeTopicName, func(message *sarama.ConsumerMessage) {
		log.Println(string(message.Value))
	})
	if err = orderStateChangeListener.Subscribe(); err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	<-context.TODO().Done()
}
