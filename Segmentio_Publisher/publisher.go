package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {

	bootstrapServers := os.Getenv("BOOTSTRAP_SERVER")
	topic := os.Getenv("TOPIC")
	partition := os.Getenv("PARTITION")

	fmt.Printf("BOOTSTRAP_SERVER: %s\nTOPIC: %s\nPARTITION: %s\n", bootstrapServers, topic, partition)

	kfClient := kafka.Client{}

	partitionInt, _ := strconv.Atoi(partition)

	conn, err := kafka.DialLeader(context.TODO(), "tcp", bootstrapServers, topic, partitionInt)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	log.Println("kafka broker connected ....")

	defer func(conn *kafka.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal("Failed to close connection to kafka", err)
		}
	}(conn)

	log.Println("Started Producing..........")

	msgCnt := 1
	for {
		time.Sleep(10 * time.Second)
		// send a batch of messages
		messages := []kafka.Message{
			{Value: []byte(fmt.Sprintf("Message - %d", msgCnt))},
		}
		_, err := conn.WriteMessages(
			messages...,
		)
		if err != nil {
			log.Println("Failed to write messages:", err)
			continue
		}
		fmt.Println("Sent Messages:")
		fmt.Printf("%s\n", string(messages[0].Value))
		msgCnt++
	}
}
