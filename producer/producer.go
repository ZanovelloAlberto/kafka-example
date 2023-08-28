package producer

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/lorenzotinfena/kafka-example/utils"
)

// The production of messages is completely asynchronous
func Main(serverSockets string, topic string) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":  serverSockets,
		"client.id":          "producer",
		"acks":               "-1",
		"enable.idempotence": "true"}) // This prevents duplication and assures order of messages
	if err != nil {
		slog.Error(err.Error())
		panic("Exiting")
	}
	defer producer.Close()

	events := producer.Events()
	i := 0
	var w bytes.Buffer
	enc := gob.NewEncoder(&w)
	latestData := utils.Data{DateTime: time.Now().UTC(), Number: i}
	enc.Encode(latestData)
	messageToSend := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          w.Bytes(),
	}

	for {
		err = producer.Produce(
			messageToSend,
			nil,
		)
		if err != nil {
			slog.Error(err.Error())
			panic("Exiting")
		}

		message, more := (<-events).(*kafka.Message)
		if !more {
			slog.Error("Cannot read events")
			panic("Exiting")
		}
		// If any error, retry sending the same message ()
		if message.TopicPartition.Error != nil {
			slog.Warn(message.TopicPartition.Error.Error())
		} else {
			slog.Info("Produced: " + fmt.Sprint(latestData.Number))

			i++

			var w bytes.Buffer
			enc := gob.NewEncoder(&w)
			latestData = utils.Data{DateTime: time.Now().UTC(), Number: i}
			enc.Encode(latestData)
			messageToSend = &kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          w.Bytes(),
			}
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(5000)))
		}
	}
}
