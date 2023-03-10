package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  "localhost:9092",
		"group.id":           "foo",
		"group.instance.id":  "1",
		"enable.auto.commit": false,
		"auto.offset.reset":  "earliest"}) // earliest, latest
	if err != nil {
		panic(err)
	}
	err = consumer.SubscribeTopics([]string{"stream"}, func(consumer *kafka.Consumer, event kafka.Event) error {
		fmt.Println("event happens: ", event.String())
		return nil
	})
	if err != nil {
		panic(err)
	}

	//justRead(consumer)
	syncCommit(consumer, 1000)

}

func justRead(consumer *kafka.Consumer) {

	for {
		message, err := consumer.ReadMessage(-1)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("consumer message content: ", string(message.Value))

	}

}

func syncCommit(consumer *kafka.Consumer, commitBatchSize int) {
	var count = 0
	for {
		message, err := consumer.ReadMessage(-1)
		if err != nil {
			fmt.Println(err)
			continue
		}

		//// actual process on data
		//fmt.Println("consumer message content: ", string(message.Value))

		count++
		if count == commitBatchSize {
			commitMessage, err := consumer.CommitMessage(message)
			if err != nil {
				fmt.Println(err)

			}
			fmt.Println("committed topic_partitions: ", commitMessage)
			count = 0
		}

	}
}
