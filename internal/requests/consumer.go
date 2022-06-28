package requests

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/amdf/conv-make-img/internal/config"
	"github.com/amdf/conv-make-img/internal/converter"
)

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	Ready chan bool
	conv  *converter.TengwarConverter
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	convAddr := config.GetConverterAddress()

	log.Println("consumer Setup - connect to ", convAddr)

	var err error
	consumer.conv, err = converter.NewTengwarConverter(convAddr)
	if err == nil {
		log.Println("success")
		// Mark the consumer as ready
		close(consumer.Ready)

		return nil
	}
	return err
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	consumer.conv.ClientConn.Close()
	log.Println("connection closed")
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message := <-claim.Messages():
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			consumer.makeImage(message.Value)
			session.MarkMessage(message, "")

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/Shopify/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}

func (consumer *Consumer) makeImage(rawMessage []byte) {
	if 0 == len(rawMessage) {
		fmt.Println("makeImage: Empty request")
		return
	}

	rq := converter.ConvertRequest{}
	err := json.Unmarshal(rawMessage, &rq)
	if err != nil {
		fmt.Println("makeImage: Unknown request")
		return
	}

	bytes, err := consumer.conv.MakeImage(context.Background(), rq)
	if err != nil {
		fmt.Println("MakeImage:", err.Error())
		return
	}

	err = converter.SaveImage(rq.ConvID, bytes)

	if err != nil {
		fmt.Println("Convert error:", err.Error())
		return
	}

	fmt.Println("Successfully convert to ", rq.ConvID, "size", len(bytes))
}
