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

	invertedCoordinates := []float64{coordinates[1], coordinates[0]} // beacause coordinates in MT are inverted (check readme)

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
				//Coordinates: coordinates,
				Coordinates: invertedCoordinates,
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
	initialOffset := u.kafkaReader.Stats().Offset

	vehicles := []*entities.Vehicle{}

L: //label used to break the for loop
	for {
		select {
		case <-stopChan:
			u.logger.Info("stopping server gracefully")
			u.kafkaReader.SetOffset(initialOffset)
			u.logger.Info("kafka reader offset resetted")
			u.kafkaReader.Close()
			return nil
		default:
			message, err := u.kafkaReader.ReadMessage(ctx)
			if err != nil {
				u.logger.Errorw("can't read message", "error", err)
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
			vehicles = append(vehicles, v)

			if u.kafkaReader.Lag() == 0 {
				u.kafkaReader.Close()
				break L
			}
		}
	}
	
	err := u.persistor.WriteVehiclesBatch(ctx, vehicles)
	if err != nil {
		u.logger.Errorw("can't write vehicles", "error", err)
		return errors.Wrap(err, "can't write vehicles")
	}
	u.logger.Infow("vehicles written", "count", len(vehicles))
	return nil
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