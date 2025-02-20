package config

import (
	"embed"
	"log"
	"strconv"

	"gopkg.in/yaml.v3"
)

var MicroConfig MicroConf

type ConsulConf struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type GrpcConf struct {
	Address            string `yaml:"address"`
	Port               int    `yaml:"port"`
	ServiceName        string `yaml:"serviceName"`
	HealthCheckSeconds int    `yaml:"healthCheckSeconds"`
}

type MicroConf struct {
	Consul ConsulConf `yaml:"consul"`
	Grpc   GrpcConf   `yaml:"grpc"`
}

func (consulConfig ConsulConf) GenConsulAddr() string {
	return consulConfig.Address + ":" + strconv.Itoa(consulConfig.Port)
}

func (this GrpcConf) GenServiceAddr() string {
	return this.Address + ":" + strconv.Itoa(this.Port)
}

//go:embed microServiceConfig.yaml
var microf embed.FS

func init() {
	configMicro()
}

func configMicro() {
	file, err := microf.ReadFile("microServiceConfig.yaml")
	if err != nil {
		log.Panic(err)
	}
	err = yaml.Unmarshal(file, &MicroConfig)
	if err != nil {
		log.Panic(err)
	}
}
