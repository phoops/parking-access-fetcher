package ngsild

import (
	"context"
	"fmt"

	"bitbucket.org/phoops/nurse/internal/core/entities"
	"github.com/philiphil/geojson"
	"github.com/phoops/ngsi-gold/client"
	"github.com/phoops/ngsi-gold/ldcontext"
	"github.com/phoops/ngsi-gold/model"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Client struct {
	logger       *zap.SugaredLogger
	baseURL      string
	ngsiLdClient *client.NgsiLdClient
}

func NewClient(logger *zap.SugaredLogger, baseURL string) (*Client, error) {
	if logger == nil {
		return nil, errors.New("all parameters must be non-nil")
	}
	logger = logger.With("component", "NGSI-LD client")
	ngsiLdClient, err := client.New(
		client.SetURL(baseURL),
	)
	if err != nil {
		return nil, errors.Wrap(err, "can't instantiate ngsi-ld client")
	}

	return &Client{
		logger,
		baseURL,
		ngsiLdClient,
	}, nil
}

func (c *Client) WriteVehiclesBatch(ctx context.Context, vehicles []*entities.Vehicle) error {
	payload := []*client.EntityWithContext{}
	for _, v := range vehicles {
		e := vehicleToBrokerEntity(v)
		defaulContext := ldcontext.LdContext{
			"https://raw.githubusercontent.com/smart-data-models/data-models/master/context.jsonld",
			"https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
		}
		payload = append(payload, &client.EntityWithContext{
			LdCtx:  &defaulContext,
			Entity: e,
		})
	}

	err := c.ngsiLdClient.BatchUpsertEntities(ctx, payload, client.UpsertSetUpdateMode)
	if err != nil {
		c.logger.Errorw("can't update entities", "err", err)
		return errors.Wrap(err, "can't update entities")
	}
	return nil
}

func vehicleToBrokerEntity(e *entities.Vehicle) *model.Entity {
	newId := fmt.Sprintf("urn:ngsi-ld:vehicle:%s", e.Id)
	location := geojson.NewPointGeometry([]float64{float64(e.Location.Value.Coordinates[1]), float64(e.Location.Value.Coordinates[0])}) // coordinates are swapped because of mobility toolkit preblems // TODO fix this when are solved

	return &model.Entity{
		ID:   newId,
		Type: "Vehicle",
		Properties: model.Properties{
			"speed": model.Property{
				Value:      e.Speed.Value,
				ObservedAt: &e.Speed.ObservedAt,
			},
			"vehicleType": model.Property{
				Value: e.VehicleType,
			},
			"description": model.Property{
				Value: e.Description,
			},
			"heading": model.Property{
				Value:      e.Heading.Value,
				ObservedAt: &e.Heading.ObservedAt,
			},
		},
		Location: &model.GeoProperty{
			Value:      location,
			ObservedAt: &e.Location.ObservedAt,
		},
	}
}
