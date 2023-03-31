package main

import (
	"context"
	"log"
	"route256/libs/kafka"
	"route256/notifications/internal/config"

	"github.com/Shopify/sarama"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal("config init", err)
	}

	brokers := config.ConfigData.Kafka.Brokers
	c, err := kafka.NewConsumer(brokers)
	if err != nil {
		log.Fatalln(err)
	}

	orderStateChangeTopicName := config.ConfigData.Kafka.Topics.OrderStateChange.Name
	orderStateChangeListener := kafka.NewListener(c, orderStateChangeTopicName, func(message *sarama.ConsumerMessage) {
		log.Println(string(message.Value))
	})
	if err = orderStateChangeListener.Subscribe(); err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	<-context.TODO().Done()
}
