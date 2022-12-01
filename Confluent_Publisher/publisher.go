package main

import (
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {

	bootstrapServers := os.Getenv("BOOTSTRAP_SERVER")
	topic := os.Getenv("TOPIC")
	totalMsgcnt := 1000

	fmt.Printf("BOOTSTRAP_SERVER: %s\nTOPIC: %s\n", bootstrapServers, topic)

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Producer %v\n", p)

	// Listen to all the client instance-level errors.
	// It's important to read these errors too otherwise the events channel will eventually fill up
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case kafka.Error:
				// Generic client instance-level errors, such as
				// broker connection failures, authentication issues, etc.
				//
				// These errors should generally be considered informational
				// as the underlying client will automatically try to
				// recover from any errors encountered, the application
				// does not need to take action on them.
				fmt.Printf("Error: %v\n", ev)
			default:
				fmt.Printf("Ignored event: %s\n", ev)
			}
		}
	}()

	msgcnt := 0
	for msgcnt < totalMsgcnt {
		value := fmt.Sprintf("Producer example, message #%d", msgcnt)

		// A delivery channel for each message sent.
		// This permits to receive delivery reports
		// separately and to handle the use case
		// of a server that has multiple concurrent
		// produce requests and needs to deliver the replies
		// to many different response channels.
		deliveryChan := make(chan kafka.Event)
		go func() {
			for e := range deliveryChan {
				switch ev := e.(type) {
				case *kafka.Message:
					// The message delivery report, indicating success or
					// permanent failure after retries have been exhausted.
					// Application level retries won't help since the client
					// is already configured to do that.
					m := ev
					if m.TopicPartition.Error != nil {
						fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
					} else {
						fmt.Printf("Delivered message< %v > to topic %s partition [%d] at offset %v\n", string(m.Value),
							*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
					}

				default:
					fmt.Printf("Ignored event: %s\n", ev)
				}
				// in this case the caller knows that this channel is used only
				// for one Produce call, so it can close it.
				close(deliveryChan)
			}
		}()

		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(value),
			Headers:        []kafka.Header{{Key: "Header", Value: []byte("Hi! consumer")}},
		}, deliveryChan)

		if err != nil {
			close(deliveryChan)
			if err.(kafka.Error).Code() == kafka.ErrQueueFull {
				// Producer queue is full, wait 1s for messages
				// to be delivered then try again.
				time.Sleep(time.Second)
				continue
			}
			fmt.Printf("Failed to produce message: %v\n", err)
		}
		msgcnt++
		time.Sleep(10 * time.Second)
	}

	// Flush and close the producer and the events channel
	for p.Flush(10000) > 0 {
		fmt.Print("Still waiting to flush outstanding messages\n", err)
	}
	p.Close()
}
