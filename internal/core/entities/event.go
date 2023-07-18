package entities

import (
	"time"

	"github.com/google/uuid"
)

type PresenceEvent struct {
	ID          uuid.UUID `json:"id"`
	Source      string    `json:"source"`
	PlateNumber string    `json:"plateNumber"`
	Country     string    `json:"country"`
	GateID      string    `json:"gateId"`
	ParkingID   string    `json:"parkingId"`
	Direction   string    `json:"direction"`
	DetectedAt  time.Time `json:"detectedAt"`
	ReceivedAt  time.Time `json:"receivedAt"`
}