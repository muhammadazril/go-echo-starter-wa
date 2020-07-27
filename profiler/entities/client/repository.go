package client

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/rimantoro/event_driven/profiler/entities/client/model"
)

type StreamRepository interface {
	StreamProduceMessages(p *kafka.Producer, msgs []string) error

	StreamForCreate(p *kafka.Producer, m *model.Client) (*model.Client, error)
	StreamForUpdate(p *kafka.Producer, m *model.Client) (*model.Client, error)
	StreamCrudProcess(c *kafka.Consumer) error
}
