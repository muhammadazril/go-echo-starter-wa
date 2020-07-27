package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"

	"github.com/rimantoro/event_driven/profiler/entities/client/model"
	"github.com/rimantoro/event_driven/profiler/shared/bootstrap"
)

func (m *clientRepository) StreamProduceMessages(p *kafka.Producer, msgs []string) error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "0.0.0.0:9092"})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	for _, word := range msgs {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &model.Topic, Partition: 3},
			Value:          []byte(word),
		}, nil)
	}

	// Wait for message deliveries before shutting down
	p.Flush(5000)
	return nil
}

func (r *clientRepository) StreamForCreate(p *kafka.Producer, m *model.Client) (*model.Client, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "0.0.0.0:9092"})
	if err != nil {
		bootstrap.App.Logger.Error("create producer failed", zap.Error(err))
		return m, err
	}

	payload := model.CrudPayload{
		ProcName: "create",
		Payload:  m,
	}

	if bPayload, err := json.Marshal(payload); err != nil {
		bootstrap.App.Logger.Error("json marshal error", zap.Error(err))
		return m, err
	} else {
		delivery_chan := make(chan kafka.Event, 100)
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &model.Topic, Partition: kafka.PartitionAny},
			Value:          bPayload},
			delivery_chan,
		)
		e := <-delivery_chan
		km := e.(*kafka.Message)
		if km.TopicPartition.Error != nil {
			bootstrap.App.Logger.Error("delivery failed", zap.Error(km.TopicPartition.Error))
		} else {
			bootstrap.App.Logger.Info("delivery success", zap.String("topic", *km.TopicPartition.Topic), zap.Int32("partition", km.TopicPartition.Partition), zap.Any("offset", km.TopicPartition.Offset), zap.Any("payload", string(km.Value)))
		}
		close(delivery_chan)
		return m, nil
	}
}

func (r *clientRepository) StreamForUpdate(k *kafka.Producer, m *model.Client) (*model.Client, error) {
	return m, errors.New("test")
}

func (r *clientRepository) StreamCrudProcess(k *kafka.Consumer) error {
	return errors.New("test")
}

func StreamConsumeCrud(c *kafka.Consumer) *kafka.Consumer {
	c.SubscribeTopics([]string{model.Topic}, nil)
	run := true
	MIN_COMMIT_COUNT := 1
	msgCount := 0
	for run == true {
		ev := c.Poll(10000)
		switch e := ev.(type) {
		case *kafka.Message:
			msgCount++
			if msgCount%MIN_COMMIT_COUNT == 0 {
				go func() {
					offsets, err := c.Commit()
					if err != nil {
						fmt.Printf("%% Error Commit on %s:\n%s\n", offsets, err.Error())
					}
				}()
			}
			fmt.Printf("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
		case kafka.PartitionEOF:
			fmt.Printf("%% Reached %v\n", e)
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			run = false
		default:
			fmt.Printf("Ignored %v\n", e)
		}
	}

	return c
}
