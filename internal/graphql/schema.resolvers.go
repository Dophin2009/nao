package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/Dophin2009/nao/pkg/data"
	"github.com/Dophin2009/nao/pkg/db"
)

func (r *mutationResolver) CreateMedia(ctx context.Context, media data.Media) (*data.Media, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	err = ds.Database.Transaction(true, func(tx db.Tx) error {
		ser := ds.MediaService
		_, err = ser.Create(&media, tx)
		if err != nil {
			return fmt.Errorf("failed to create Media: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &media, nil
}

func (r *queryResolver) MediaByID(ctx context.Context, id int) (*data.Media, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var md *data.Media
	err = ds.Database.Transaction(false, func(tx db.Tx) error {
		ser := ds.MediaService
		md, err = ser.GetByID(id, tx)
		if err != nil {
			return fmt.Errorf("failed to get Media by id %d: %w", id, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return md, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
