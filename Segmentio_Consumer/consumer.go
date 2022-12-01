package main

// consumer_example implements a consumer using the non-channel Poll() API
// to retrieve messages and events.

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"strings"
)

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 1,    // 1B
		MaxBytes: 10e6, // 10MB
	})
}

func main() {

	bootstrapServers := os.Getenv("BOOTSTRAP_SERVER")
	topic := os.Getenv("TOPIC")
	partition := os.Getenv("PARTITION")

	fmt.Printf("BOOTSTRAP_SERVER: %s\nTOPIC: %s\nPARTITION: %v\n", bootstrapServers, topic, partition)

	reader := getKafkaReader(bootstrapServers, topic, "one")

	defer reader.Close()

	fmt.Println("start consuming ... !!")
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Printf("Received Message: %s <at topic:%v partition:%v offset:%v>\n", string(m.Value), m.Topic, m.Partition, m.Offset)
	}
}
