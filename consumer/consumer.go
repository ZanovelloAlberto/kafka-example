package consumer

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log/slog"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/lorenzotinfena/goji/collections"
	"github.com/lorenzotinfena/kafka-example/utils"
)

type Consumer struct {
	client *kafka.Consumer
	topic  string
}

// Read all past messages, plus one future message, if no message are read in 10 seconds, it will return
func (con *Consumer) ReadAllUntilNow() []string {
	err := con.client.Subscribe(con.topic, nil)
	if err != nil {
		slog.Error(err.Error())
		panic("Exiting")
	}

	results := collections.NewSingleLinkedList[string](nil)
	now := time.Now()
	for {
		message, err := con.client.ReadMessage(time.Second * 10)
		if err != nil {
			switch e := err.(type) {
			case kafka.Error:
				if e.Code() == kafka.ErrTimedOut {
					return results.ToSlice()
				}
			default:
			}
			slog.Warn(err.Error())
			continue
		}

		var r bytes.Buffer
		r.Write(message.Value)
		enc := gob.NewDecoder(&r)
		var data utils.Data
		enc.Decode(&data)

		results.InsertLast(data.String())
		slog.Info("Consumed: " + fmt.Sprint(data.Number))
		if message.Timestamp.After(now) {
			return results.ToSlice()
		}
	}
}
func (con *Consumer) Close() {
	con.client.Close()
}
func NewConsumer(serverSockets string, topic string) *Consumer {
	client, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": serverSockets,
		"group.id":          "consumer",
		"auto.offset.reset": "smallest"})
	if err != nil {
		slog.Error(err.Error())
		panic("Exiting")
	}
	return &Consumer{client: client, topic: topic}
}
