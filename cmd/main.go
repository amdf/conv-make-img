package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/amdf/conv-make-img/internal/config"
	"github.com/amdf/conv-make-img/internal/requests"
)

// Sarama configuration options
var (
	verbose = false
)

func main() {
	if nil != config.Load() {
		log.Fatal("unable to load configs/config.toml")
	}

	keepRunning := true
	log.Println("Starting a new Sarama consumer")

	if config.Get().Consumer.Verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	version, err := sarama.ParseKafkaVersion("2.1.1")
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	kafkacfg := sarama.NewConfig()
	kafkacfg.Version = version
	kafkacfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	kafkacfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	/**
	 * Setup a new Sarama consumer group
	 */
	consumer := requests.Consumer{
		Ready: make(chan bool),
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(config.Get().Consumer.Brokers, config.Get().Consumer.Group, kafkacfg)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, []string{config.Get().Consumer.Topic}, &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.Ready = make(chan bool)
		}
	}()

	<-consumer.Ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		}
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}
