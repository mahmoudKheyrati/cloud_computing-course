package main

import (
	"flag"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

var (
	groupId         = ""
	groupInstanceId = ""
	commitMode      = ""
	batchSize       = 1
)

func main() {
	flag.StringVar(&groupId, "groupId", "foo", "groupId of the kafka consumer")
	flag.StringVar(&groupInstanceId, "groupInstanceId", "1", "groupInstanceId of the kafka consumer")
	flag.StringVar(&commitMode, "commitMode", "justRead", "commit mode of the kafka consumer. should be one of these values: justRead, syncCommit, asyncCommit")
	flag.IntVar(&batchSize, "batchSize", 1, "batchSize to commit")
	flag.Parse()
	log.Println("start kafkaConsumer with {groupId=", groupId, " and groupInstanceId=", groupInstanceId, "commitMode=", commitMode, " batchSize=", batchSize, "}")

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  "localhost:9092",
		"group.id":           groupId,
		"group.instance.id":  groupInstanceId,
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

	switch commitMode {
	case "justRead":
		justRead(consumer)
	case "syncCommit":
		syncCommit(consumer, batchSize)
	case "asyncCommit":
		AsyncCommit(consumer, batchSize)
	default:
		log.Println("commitMode should be on of these values: justRead, syncCommit, asyncCommit")
	}

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
		fmt.Println("consumer message content: ", string(message.Value))

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

func AsyncCommit(consumer *kafka.Consumer, commitBatchSize int) {
	var count = 0
	for {
		message, err := consumer.ReadMessage(-1)
		if err != nil {
			fmt.Println(err)
			continue
		}

		//// actual process on data
		fmt.Println("consumer message content: ", string(message.Value))

		count++
		if count == commitBatchSize {
			go func() {
				commitMessage, err := consumer.CommitMessage(message)
				if err != nil {
					fmt.Println(err)

				}
				fmt.Println("committed topic_partitions: ", commitMessage)
				count = 0
			}()
		}

	}
}
