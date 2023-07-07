package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"bitbucket.org/phoops/nurse/internal/core/entities"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type VehiclePersistor interface {
	WriteVehiclesBatch(ctx context.Context, vs []*entities.Vehicle) error
}

type SyncVehicles struct {
	logger              *zap.SugaredLogger
	kafkaReader         *kafka.Reader
	persistor           VehiclePersistor
	defaultVehicleSpeed int
	kafkaTopic          string
}

func NewSyncVehicle(
	logger *zap.SugaredLogger,
	kafkaReader *kafka.Reader,
	persistor VehiclePersistor,
	defaultVehicleSpeed int,
	kafkaTopic string,
) (*SyncVehicles, error) {
	if logger == nil || persistor == nil || kafkaReader == nil {
		return nil, errors.New("all parameters must be non-nil")
	}
	logger = logger.With("usecase", "SyncVehicles")

	return &SyncVehicles{
		logger,
		kafkaReader,
		persistor,
		defaultVehicleSpeed,
		kafkaTopic,
	}, nil
}

func (u *SyncVehicles) presenceEvent2Vehicle(pe entities.PresenceEvent) (*entities.Vehicle, error) {

	var coordinates []float64
	switch pe.ParkingID {
	case "atam-off-street-parking-cadorna":
		coordinates = []float64{43.465313, 11.872549}
	case "atam-off-street-parking-san-donato":
		coordinates = []float64{43.462014, 11.864127}
	case "atam-off-street-parking-baldaccio":
		coordinates = []float64{43.465313, 11.872549}
	case "atam-off-street-parking-mecenate":
		coordinates = []float64{43.455705, 11.880767}
	default:
		u.logger.Errorw("parking ID not found", "parkingID", pe.ParkingID)
		return nil, errors.New("parking ID not found")
	}

	return &entities.Vehicle{
		Id:          pe.ID.String(),
		Type:        "Vehicle",
		VehicleType: "Car",
		Speed: entities.Speed{
			Value:      u.defaultVehicleSpeed,
			ObservedAt: pe.DetectedAt,
		},
		Location: entities.Location{
			Value: entities.Point{
				Coordinates: coordinates,
			},
			ObservedAt: pe.DetectedAt,
		},
		Description: fmt.Sprintf("Parking: %s, Gate: %s", pe.ParkingID, pe.GateID),
		Heading: entities.Heading{
			Value:      180, //default value. Not used
			ObservedAt: pe.DetectedAt,
		},
	}, nil
}

func (u *SyncVehicles) Execute(ctx context.Context) error {
	u.logger.Info("running vehicles synchronization")

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	

	for {
		select {
		case <-stopChan:
			u.logger.Info("stopping server gracefully")
			u.kafkaReader.Close()
			return nil
		default:
			message, err := u.kafkaReader.ReadMessage(ctx)
			if err != nil {
				u.logger.Errorw("can't read vehicle message", "error", err)
				return errors.Wrap(err, "can't read vehicle message")
			}

			var presenceEvent entities.PresenceEvent
			err = json.Unmarshal(message.Value, &presenceEvent)
			if err != nil {
				u.logger.Errorw("can't unmarshal vehicle message", "error", err)
				return errors.Wrap(err, "can't unmarshal vehicle message")
			}

			u.logger.Debugw("message received", "message", presenceEvent)
			v, err := u.presenceEvent2Vehicle(presenceEvent)
			if err != nil {
				u.logger.Errorw("can't convert presence event to vehicle", "error", err)
				if u.kafkaReader.Lag() > 0 {
					continue
				}
			}

		
			err = u.persistor.WriteVehiclesBatch(ctx, []*entities.Vehicle{v})
			if err != nil {
				u.logger.Errorw("can't write vehicle", "error", err)
				return errors.Wrap(err, "can't write vehicle")
			}
			u.logger.Infow("vehicle written. Note: coordinates will be inverted (check readme) ", "vehicle", v)
		}
			
			
	}

}