package config

import (
	"embed"
	"log"

	"gopkg.in/yaml.v3"
)

var KafkaConfig KafkaConf

type KafkaConf struct {
	VideoPub     ConsumerConf `yaml:"videoPub"`
	VideoTrscode ProducerConf `yaml:"videoTrscode"`
}

type ConsumerConf struct {
	Topic            string   `yaml:"topic"`
	Group            string   `yaml:"group"`
	Brokers          []string `yaml:"brokers"`
	MinBytes         int      `yaml:"minBytes"`
	MaxBytes         int      `yaml:"maxBytes"`
	ReadBatchTimeout int      `yaml:"readBatchTimeOut"`
}

type ProducerConf struct {
	Topic        string   `yaml:"topic"`
	Brokers      []string `yaml:"brokers"`
	WriteTimtout int      `yaml:"writeTimeout"`
	Acks         int      `yaml:"acks"`
	BatchSize    int      `yaml:"batchSize"`
}

//go:embed kafkaConfig.yaml
var kafkaFs embed.FS

func init() {
	configKafka()
}

func configKafka() {
	file, err := kafkaFs.ReadFile("kafkaConfig.yaml")
	if err != nil {
		log.Panic("start up to config kafka failure")
	}
	err = yaml.Unmarshal(file, &KafkaConfig)
	if err != nil {
		log.Panic("bind config file to kafkaConfig failure")
	}

}
