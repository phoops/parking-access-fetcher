package usecase

import (
	"context"
	"fmt"

	//"time"

	"bitbucket.org/phoops/nurse/internal/core/entities"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// type vsFetcher interface {
// 	Getvs(ctx context.Context) ([]*entities.v, error)
// }

type VehiclePersistor interface {
	WriteVehiclesBatch(ctx context.Context, vs []*entities.Vehicle) error
}

type SyncVehicles struct {
	logger      *zap.SugaredLogger
	kafkaReader *kafka.Reader
	persistor   VehiclePersistor
}

func NewSyncVehicle(
	logger *zap.SugaredLogger,
	kafkaReader *kafka.Reader,
	persistor VehiclePersistor,
) (*SyncVehicles, error) {
	if logger == nil || persistor == nil || kafkaReader == nil {
		return nil, errors.New("all parameters must be non-nil")
	}
	logger = logger.With("usecase", "SyncVehicles")

	return &SyncVehicles{
		logger,
		kafkaReader,
		persistor,
	}, nil
}

func (u *SyncVehicles) Execute(ctx context.Context) error {
	u.logger.Info("running vehicles synchronization")

	for {
		message, err := u.kafkaReader.ReadMessage(ctx)
		if err != nil {
			u.logger.Errorw("can't read message", "error", err)
			return errors.Wrap(err, "can't read stops")
		}
		u.logger.Debugw("messages received", "message", message)
	}

	// +++++++++++++++ create mockup vehicle data for testing +++++++++++++++
	// vehicles := []*entities.Vehicle{}
	// for i := 1; i <= 100; i++ {
	// 	v := &entities.Vehicle{
	// 		Id:          fmt.Sprintf("%s%03d", "urn:ngsi-ld:Vehicle:", i),
	// 		Type:        "Vehicle",
	// 		VehicleType: "car",
	// 		Description: "camera 1",
	// 		Speed: entities.Speed{
	// 			Value:      50,
	// 			ObservedAt: time.Now(),
	// 		},
	// 		Location: entities.Location{
	// 			Value: entities.Point{
	// 				Coordinates: []float64{43.459137, 11.861667},
	// 			},
	// 			ObservedAt: time.Now(),
	// 		},
	// 		Heading: entities.Heading{
	// 			Value:      180,
	// 			ObservedAt: time.Now(),
	// 		},
	// 	}
	// 	vehicles = append(vehicles, v)
	// }

	// +++++++++++++++ write vehicles on broker +++++++++++++++ //TODO uncomment
	// err := u.persistor.WriteVehiclesBatch(ctx, vehicles)
	// if err != nil {
	// 	u.logger.Errorw("can't write vehicles", "error", err)
	// 	return errors.Wrap(err, "can't write vehicles")
	// }
	// u.logger.Infow("vehicles written", "size", len(vehicles))
	// return nil
}
