package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
)

var brokers = []string{
	"localhost:9095",
	"localhost:9096",
	"localhost:9097",
}

var flagNumPart = flag.Int("n", 1, "number of partition (from 1 to n)")

func main() {
	flag.Parse()
	if *flagNumPart <= 0 {
		log.Fatalln("number of partition should start from 1")
	}

	config := sarama.NewConfig()
	config.ClientID = "conv-make-img"
	config.Consumer.Return.Errors = true

	// Create new consumer
	convConsumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := convConsumer.Close(); err != nil {
			panic(err)
		}
	}()

	//topics, _ := convConsumer.Topics()
	ctx, cancel := context.WithCancel(context.Background())

	numPart := *flagNumPart - 1
	messages, errors, err := startConsume(ctx, "convert_requests", numPart, convConsumer)
	if err != nil {
		log.Fatalln(err.Error())
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Count how many message processed
	msgCount := 0

	// Get signal for finish
	done := make(chan struct{})
	go func() {
		for {
			select {
			case msg := <-messages:
				msgCount++
				fmt.Println("Received messages", string(msg.Key), string(msg.Value))
			case consumerError := <-errors:
				msgCount++
				fmt.Println("Received consumerError ", string(consumerError.Topic), string(consumerError.Partition), consumerError.Err)
				done <- struct{}{}
			case <-signals:
				fmt.Println("Interrupt is detected")
				done <- struct{}{}
			}
		}
	}()

	<-done
	cancel()
	fmt.Println("Processed", msgCount, "messages")

}

func startConsume(ctx context.Context, topic string, numPart int, consumer sarama.Consumer) (msgch chan *sarama.ConsumerMessage, errch chan *sarama.ConsumerError, err error) {
	msgch = make(chan *sarama.ConsumerMessage)
	errch = make(chan *sarama.ConsumerError)

	partitions, _ := consumer.Partitions(topic)

	if numPart >= len(partitions) {
		err = errors.New("nonexistent partition")
		return
	}

	partConsumer, err := consumer.ConsumePartition(topic, partitions[numPart], sarama.OffsetOldest)
	if nil != err {
		fmt.Printf("Topic %s Partition: %d of %d", topic, numPart, len(partitions))
		panic(err)
	}

	fmt.Println(" Start consuming topic ", topic)
	go func(ctx context.Context, topic string, consumer sarama.PartitionConsumer) {
		consuming := true
		for consuming {
			select {
			case <-ctx.Done():
				fmt.Println("canceled by context")
				consuming = false
			case consumerError := <-consumer.Errors():
				errch <- consumerError
				fmt.Println("consumerError: ", consumerError.Err)

			case msg := <-consumer.Messages():
				msgch <- msg
				fmt.Println("Got message on topic ", topic, msg.Value)
			}
		}

		fmt.Println("end consuming")
	}(ctx, topic, partConsumer)

	return
}
