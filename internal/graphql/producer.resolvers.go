package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/Dophin2009/nao/pkg/db"
	"github.com/Dophin2009/nao/pkg/models"
)

func (r *producerResolver) Titles(ctx context.Context, obj *models.Producer, first *int, skip *int) ([]*models.Title, error) {
	return sliceTitles(obj.Titles, first, skip), nil
}

func (r *producerResolver) Media(ctx context.Context, obj *models.Producer, first *int, skip *int) ([]*models.MediaProducer, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var list []*models.MediaProducer
	err = ds.Database.Transaction(false, func(tx db.Tx) error {
		ser := ds.MediaProducerService
		list, err = ser.GetByProducer(obj.Meta.ID, first, skip, tx)
		if err != nil {
			return fmt.Errorf(
				"failed to get MediaProducers by Producer id %d: %w", obj.Meta.ID, err)
		}
		return nil
	})

	return list, nil
}

// Producer returns ProducerResolver implementation.
func (r *Resolver) Producer() ProducerResolver { return &producerResolver{r} }

type producerResolver struct{ *Resolver }
