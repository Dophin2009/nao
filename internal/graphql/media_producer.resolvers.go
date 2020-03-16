package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/Dophin2009/nao/pkg/data/models"
	"github.com/Dophin2009/nao/pkg/db"
)

func (r *mediaProducerResolver) Media(ctx context.Context, obj *models.MediaProducer) (*models.Media, error) {
	return resolveMediaByID(ctx, obj.MediaID)
}

func (r *mediaProducerResolver) Producer(ctx context.Context, obj *models.MediaProducer) (*models.Producer, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var p *models.Producer
	err = ds.Database.Transaction(false, func(tx db.Tx) error {
		ser := ds.ProducerService
		p, err = ser.GetByID(obj.ProducerID, tx)
		if err != nil {
			return fmt.Errorf("failed to get Producer by id %d: %w", obj.ProducerID, err)
		}
		return nil
	})

	return p, nil
}

// MediaProducer returns MediaProducerResolver implementation.
func (r *Resolver) MediaProducer() MediaProducerResolver { return &mediaProducerResolver{r} }

type mediaProducerResolver struct{ *Resolver }
