package stream

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// consumer for CRUD process
func ConsumeTestMessages(c *kafka.Consumer, topics []string) *kafka.Consumer {
	c.SubscribeTopics(topics, nil)
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
