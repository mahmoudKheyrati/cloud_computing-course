package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
	"sync"
)

func main() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"acks":              "all"})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}
	produceMessage(p, 1)

	// produce 10000 messages concurrently
	var wg sync.WaitGroup
	for i := 0; i < 10_000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			produceMessage(p, i)
		}()
	}
	wg.Wait()
	fmt.Println("all messages published to kafka")

}

func produceMessage(p *kafka.Producer, i int) {
	topic := "stream"
	deliveryChannel := make(chan kafka.Event)

	err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(fmt.Sprintf("value %d", i)),
	}, deliveryChannel)
	if err != nil {
		fmt.Println("produce message error: ", err)
	}
	<-deliveryChannel
	//fmt.Println("message published to kafka")
}
