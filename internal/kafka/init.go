package kafka

import (
	"log"
	msgConsumer "muses-engine/internal/kafka/consumer"
	msgProducer "muses-engine/internal/kafka/producer"
)

func InitKafka() {
	log.Println("begin init kafka client")
	msgConsumer.StartConsumer()
	msgProducer.StartProducer()

}
