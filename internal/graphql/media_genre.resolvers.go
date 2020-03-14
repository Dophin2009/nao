package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.com/Dophin2009/nao/pkg/data"
	"gitlab.com/Dophin2009/nao/pkg/db"
)

func (r *mediaGenreResolver) Media(ctx context.Context, obj *data.MediaGenre) (*data.Media, error) {
	return resolveMediaByID(ctx, obj.MediaID)
}

func (r *mediaGenreResolver) Genre(ctx context.Context, obj *data.MediaGenre) (*data.Genre, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var g *data.Genre
	err = ds.Database.Transaction(false, func(tx db.Tx) error {
		ser := ds.GenreService
		g, err = ser.GetByID(obj.GenreID, tx)
		if err != nil {
			return fmt.Errorf("failed to get Genre by id %d: %w", obj.GenreID, err)
		}
		return nil
	})

	return g, nil
}

// MediaGenre returns MediaGenreResolver implementation.
func (r *Resolver) MediaGenre() MediaGenreResolver { return &mediaGenreResolver{r} }

type mediaGenreResolver struct{ *Resolver }
