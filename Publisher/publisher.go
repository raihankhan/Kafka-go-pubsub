package main

import (
	"fmt"
	kafkago "github.com/Shopify/sarama"
	"log"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := produce(&wg)
		if err != nil {
			log.Println(err.Error(), "Failed to produce continuous messages")
		}
	}()
	wg.Wait()
}

func produce(wg *sync.WaitGroup) error {
	defer wg.Done()
	clientConfig := kafkago.NewConfig()
	clientConfig.Producer.Partitioner = kafkago.NewRandomPartitioner
	clientConfig.Producer.RequiredAcks = kafkago.WaitForAll
	clientConfig.Producer.Return.Successes = true
	client, err := kafkago.NewClient(
		[]string{"kafka-system-broker.demo.svc.cluster.local:9092"},
		clientConfig,
	)

	if err != nil {
		log.Println(err.Error(), "Failed to create kafka client")
		return err
	}

	producer, err := kafkago.NewSyncProducerFromClient(client)
	if err != nil {
		log.Println(err.Error(), "Failed to create new producer")
		return err
	}
	defer func(producer kafkago.SyncProducer) {
		err := producer.Close()
		if err != nil {
			log.Println(err.Error(), "Failed to close producer")
			return
		}
	}(producer)

	err = client.RefreshMetadata()
	if err != nil {
		log.Println(err, "Failed to refresh metadata")
	}

	i := 0

	for {
		i++
		message := kafkago.ProducerMessage{
			Topic:     "test",
			Value:     kafkago.StringEncoder(fmt.Sprintf("Message %v produced", i)),
			Partition: 0,
		}

		leader, err := client.Leader(message.Topic, message.Partition)
		if err != nil {
			log.Println("Failed to get leader for topic: ", message.Topic)
			continue
		}
		req := kafkago.ProduceRequest{
			RequiredAcks: 2,
		}
		res, err := leader.Produce(&req)

		partition, offset, err := producer.SendMessage(&message)
		if err != nil {
			log.Println(err.Error(), "Failed to sent message")
			continue
		}

		log.Println(fmt.Sprintf("Sent message %v to partition %v at offset %v", i, partition, offset))
		// sleep for a second
		time.Sleep(time.Second * 5)
	}

}
