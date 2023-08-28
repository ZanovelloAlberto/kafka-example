package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"github.com/lorenzotinfena/kafka-example/consumer"
	"github.com/lorenzotinfena/kafka-example/producer"
)

var topic = "veryImportantTopic"

func main() {
	serverSockets := "localhost:9092"
	admin, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": serverSockets})
	if err != nil {
		slog.Error(err.Error())
		panic("Exiting")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if ctx.Err() != nil {
		panic("Exiting")
	}
	_, err = admin.DeleteTopics(ctx, []string{topic})
	if err != nil {
		slog.Error(err.Error())
		panic("Exiting")
	}
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if ctx.Err() != nil {
		panic("Exiting")
	}
	_, err = admin.CreateTopics(ctx, []kafka.TopicSpecification{{Topic: topic, NumPartitions: 1, ReplicationFactor: 1}})
	if err != nil {
		slog.Error(err.Error())
		panic("Exiting")
	}
	go producer.Main(serverSockets, topic)
	con := consumer.NewConsumer(serverSockets, topic)
	defer con.Close()
	eng := gin.Default()
	eng.GET("/read", func(ctx *gin.Context) {
		ctx.JSON(200, con.ReadAllUntilNow())
	})
	eng.GET("/createtopic/:topicName", func(ctx *gin.Context) {
		topicName := ctx.Param("topicName")
		res, err := admin.CreateTopics(context.Background(), []kafka.TopicSpecification{{Topic: topicName, NumPartitions: 1, ReplicationFactor: 1}})
		// Trying creating a topic with an invalid name it seems is not considered an error
		if err != nil {
			ctx.Abort()
			return
		}
		ctx.JSON(200, res)
	})
	eng.Run(":8080")
}

/*
Build a Go based webserver with two goroutines such that, one goroutine pushes data to a
kafka producer, and another goroutine reads the data from kafka consumer per request.
● Introduce some random delay in the message push to kafka, have an integer that
increases for every next message, and a date time time.
●Kafka message data format : (date-UTC format, integer)
●Create an endpoint such that a kafka topic could be created whenever a request is
made. The request should contain the name of the kafka-topic.
*/
