package config

import (
	"fmt"
	"log"
	
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type NurseConfig struct {
	BrokerURL         string  `required:"true" split_words:"true"`
	KafkaURL 		  string  `required:"true" split_words:"true"`
	KafkaTopic 		  string  `required:"true" split_words:"true"`
}

func (s NurseConfig) String() string {
	return fmt.Sprintf(`
		BrokerURL: %s,
		KafkaURL: %s,
		KafkaTopic: %s,
		`,
		s.BrokerURL,
		s.KafkaURL,
		s.KafkaTopic,
	)
}

func LoadNurseConfig() (*NurseConfig, error) {
	err := godotenv.Load(".env.example") //TODO: change to .env
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