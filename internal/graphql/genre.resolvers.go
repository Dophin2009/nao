package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.com/Dophin2009/nao/pkg/data"
	"gitlab.com/Dophin2009/nao/pkg/db"
)

func (r *genreResolver) Names(ctx context.Context, obj *data.Genre, first *int, skip *int) ([]*data.Title, error) {
	return sliceTitles(obj.Names, first, skip), nil
}

func (r *genreResolver) Descriptions(ctx context.Context, obj *data.Genre, first *int, skip *int) ([]*data.Title, error) {
	return sliceTitles(obj.Descriptions, first, skip), nil
}

func (r *genreResolver) Media(ctx context.Context, obj *data.Genre, first *int, skip *int) ([]*data.MediaGenre, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var list []*data.MediaGenre
	err = ds.Database.Transaction(false, func(tx db.Tx) error {
		ser := ds.MediaGenreService
		list, err = ser.GetByGenre(obj.Meta.ID, first, skip, tx)
		if err != nil {
			return fmt.Errorf("failed to get MediaGenres by Genre id %d: %w",
				obj.Meta.ID, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}

// Genre returns GenreResolver implementation.
func (r *Resolver) Genre() GenreResolver { return &genreResolver{r} }

type genreResolver struct{ *Resolver }
