package entities

import "time"

// from FIWARE data-model https://github.com/FIWARE/data-models/blob/master/specs/Transportation/Vehicle/Vehicle/doc/spec.md4
type Vehicle struct {
    Id                    string `json:"id"`
    Type                  string `json:"type"`
    Speed                 Speed  `json:"speed"`
    Location              Location `json:"location"`
    VehicleType           string `json:"vehicleType"`
    Description           string `json:"description"`
    Heading               Heading `json:"heading"`
}

type Speed struct {
    Value      int    `json:"value"`
    ObservedAt time.Time   `json:"observedAt"`
}

type Location struct {
    Value        Point       `json:"coordinates"`
    ObservedAt   time.Time   `json:"observedAt"`
}

type Point struct {
    Type        string    `json:"type"`
    Coordinates []float64 `json:"coordinates"`
}

type Heading struct {
    Value      int    `json:"value"`
    ObservedAt time.Time `json:"observedAt"`
}

