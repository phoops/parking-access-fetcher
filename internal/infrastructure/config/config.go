package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type NurseConfig struct {
	BrokerURL           string `required:"true" split_words:"true"`
	KafkaURL            string `required:"true" split_words:"true"`
	KafkaTopic          string `required:"true" split_words:"true"`
	KafkaConsumerGroup  string `required:"true" split_words:"true"`
	DefaultVehicleSpeed int    `required:"true" split_words:"true"`
}

func (s NurseConfig) String() string {
	return fmt.Sprintf(`
		BrokerURL: %s,
		KafkaURL: %s,
		KafkaTopic: %s,
		KafkaConsumerGroup: %s,
		DefaultVehicleSpeed: %d,
		`,
		s.BrokerURL,
		s.KafkaURL,
		s.KafkaTopic,
		s.KafkaConsumerGroup,
		s.DefaultVehicleSpeed,
	)
}

func LoadNurseConfig() (*NurseConfig, error) {
	//err := godotenv.Load(".env.example")
	err := godotenv.Load()

	if err != nil {
		log.Printf("could not load configuration from .env file: %v", err)
	}
	var c NurseConfig
	err = envconfig.Process("", &c)
	if err != nil {
		return nil, err
	}
	log.Printf("Loaded configuration:%+s", c)
	return &c, nil
}
