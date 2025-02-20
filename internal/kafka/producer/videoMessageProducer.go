package producer

import (
	"context"
	"log"
	"muses-engine/config"
	"time"

	"github.com/segmentio/kafka-go"
)

var producer *kafka.Writer

func StartProducer() {
	producer = &kafka.Writer{
		Balancer:     &kafka.Hash{},
		Addr:         kafka.TCP(config.KafkaConfig.VideoTrscode.Brokers...),
		Topic:        config.KafkaConfig.VideoTrscode.Topic,
		WriteTimeout: time.Duration(config.KafkaConfig.VideoTrscode.WriteTimtout) * time.Second,
		RequiredAcks: kafka.RequiredAcks(config.KafkaConfig.VideoTrscode.Acks),
		BatchSize:    config.KafkaConfig.VideoTrscode.BatchSize,
	}
}

func SendTrsCodeFinishMsg(message []byte, key string) error {
	ctx := context.Background()
	var msg kafka.Message
	if key == "" {
		msg = kafka.Message{
			Value: message,
		}
	} else {
		msg = kafka.Message{
			Key:   []byte(key),
			Value: []byte(message),
		}
	}
	log.Println("send video transcode finished message to server ", string(message))
	err := producer.WriteMessages(ctx, msg)
	if err != nil {
		log.Println("send video transcode message failure , error is ", err)
		return err
	}
	return nil
}
