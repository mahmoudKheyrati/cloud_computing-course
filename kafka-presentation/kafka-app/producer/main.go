package main

import (
	"flag"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"os"
	"sync"
	"time"
)

var (
	messageCount = 1
)

func main() {
	flag.IntVar(&messageCount, "messageCount", 1, "count of the messages that you want to publish")
	flag.Parse()
	log.Println("start kafka producer with messageCount=", messageCount)
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"acks":              "all"})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}
	// produce 30000 messages concurrently
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			produceMessage(p, i)
		}(i)
	}
	wg.Wait()
	fmt.Println("all messages published to kafka in ", time.Since(start))

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
