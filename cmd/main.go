package main

import (
	"log"
	"muses-engine/config"
	"muses-engine/internal/app/logic/mediaStream"
	musesKafka "muses-engine/internal/kafka"
	"muses-engine/internal/micro"
	"muses-engine/internal/router"
)

func main() {
	startServer()
}

func startServer() {
	micro.RegistConsul()
	mediaStream.StartTurn()
	musesKafka.InitKafka()

	r := router.NewRouter()
	log.Println("server run in port", config.Server.HttpConfig.HttpPort)
	err := r.Run(":" + config.Server.HttpConfig.HttpPort)
	if err != nil {
		log.Println("err while run gin server ", err)
	}
}
