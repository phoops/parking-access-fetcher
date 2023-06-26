package main

import (
	"context"

	"bitbucket.org/phoops/nurse/internal/core/usecase"
	"bitbucket.org/phoops/nurse/internal/infrastructure/config"
	ngsild "bitbucket.org/phoops/nurse/internal/infrastructure/ngsi-ld"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"github.com/segmentio/kafka-go"

)

func main() {

	// Logger
	sourLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger := sourLogger.Sugar()
	
	// Configuration
	conf, err := config.LoadNurseConfig() 
	if err != nil {
		errMsg := errors.Wrap(err, "cannot read configuration").Error()
		logger.Fatal(errMsg)
	}

	// Kafka Client
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
        Brokers: []string{conf.KafkaURL},
        Topic:   conf.KafkaTopic,
    })
	
	// Context Broker Client
	contextBrokerClient, err := ngsild.NewClient(
		logger,
		conf.BrokerURL,
	)
	if err != nil {
		errMsg := errors.Wrap(err, "cannot instantiate context broker client").Error()
		logger.Fatal(errMsg)
	}

	// Sync Vehicle instance
	syncVehicle, err := usecase.NewSyncVehicle(
		logger,
		kafkaReader,
		contextBrokerClient,
	)
	if err != nil {
		errMsg := errors.Wrap(err, "cannot instantiate Nurse!").Error()
		logger.Fatal(errMsg)
	}
	logger.Infof("initialized Nurse!")

	// Execute
	err = syncVehicle.Execute(
		context.Background(),
	)
	if err != nil {
		errMsg := errors.Wrap(err, "cannot sync vehicles on context broker").Error()
		logger.Fatal(errMsg)
	}
}
