package bootstrap

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

func InitKafkaProducer(cfg kafka.ConfigMap) *kafka.Producer {
	p, err := kafka.NewProducer(&cfg)
	if err != nil {
		App.Logger.Error("create producer failed", zap.Error(err))
	}
	return p
}

func InitKafkaConsumer(cfg kafka.ConfigMap) *kafka.Consumer {
	c, err := kafka.NewConsumer(&cfg)
	if err != nil {
		c.Close()
		App.Logger.Error("consumer start failed", zap.Error(err))
	}
	return c
}
