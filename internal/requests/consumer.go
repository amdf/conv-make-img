package requests

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/amdf/conv-make-img/internal/config"
	"github.com/amdf/conv-make-img/internal/converter"
	opentracing "github.com/opentracing/opentracing-go"
	tags "github.com/opentracing/opentracing-go/ext"
)

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	Ready  chan bool
	conv   *converter.TengwarConverter
	Tracer opentracing.Tracer
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	convAddr := config.GetConverterAddress()

	log.Println("consumer Setup - connect to ", convAddr)

	var err error
	consumer.conv, err = converter.NewTengwarConverter(
		//consumer.Tracer,
		opentracing.GlobalTracer(),
		convAddr)
	if err != nil {
		return err
	}

	log.Println("success")
	// Mark the consumer as ready
	close(consumer.Ready)

	return nil
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
			tr := opentracing.GlobalTracer().StartSpan("Consume")
			tags.SpanKindConsumer.Set(tr)

			ctx := opentracing.ContextWithSpan(context.Background(), tr)
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			consumer.makeImage(ctx, message.Value)

			session.MarkMessage(message, "")

			tr.Finish()

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/Shopify/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}

func (consumer *Consumer) makeImage(ctx context.Context, rawMessage []byte) {
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

	bytes, err := consumer.conv.MakeImage(ctx, rq)
	if err != nil {
		fmt.Println("MakeImage:", err.Error())
		return
	}

	err = converter.SaveImage(ctx, rq.ConvID, bytes)

	if err != nil {
		fmt.Println("Convert error:", err.Error())
		return
	}

	fmt.Println("Successfully convert to ", rq.ConvID, "size", len(bytes))
}
