package consumer

import (
	"context"
	"encoding/json"
	"log"
	"muses-engine/config"
	"time"

	"github.com/segmentio/kafka-go"
	// "muses-engine/internal/model/message"
	videoTranscode "muses-engine/internal/app/logic/video"
	kafkaMsg "muses-engine/internal/model/message"
)

func StartConsumer() {
	go func() {
		consumeVideoState()
	}()
}

func consumeVideoState() {
	log.Println("start up kafka consumer to consume video pub message")
	// topic := "video_pub"
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:          config.KafkaConfig.VideoPub.Brokers,
		Topic:            config.KafkaConfig.VideoPub.Topic,
		GroupID:          config.KafkaConfig.VideoPub.Group,
		MinBytes:         config.KafkaConfig.VideoPub.MinBytes,
		MaxBytes:         config.KafkaConfig.VideoPub.MaxBytes,
		ReadBatchTimeout: time.Duration(config.KafkaConfig.VideoPub.ReadBatchTimeout) * time.Second,
	})

	defer reader.Close()
	ctx := context.Background()
	for {
		message, err := reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("consume message failure %V", err)
			//这里就不要break了，不然稍微抖动就会退出，但是如果一直阻塞在这里也有问题
		}
		log.Printf("receive video pub message %s ", string(message.Value))

		if err := reader.CommitMessages(ctx, message); err != nil {
			log.Printf("commit message failure err is %v ", err)
		}
		var videoPubMsg kafkaMsg.VideoPubMsg
		err = json.Unmarshal(message.Value, &videoPubMsg)
		if err != nil {
			log.Printf("unmarshal video pub message to entity failure %v", err)
		}
		videoTranscode.TranscodeVideoFile(videoPubMsg)

	}
}
