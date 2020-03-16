package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/Dophin2009/nao/pkg/data/models"
	"github.com/Dophin2009/nao/pkg/db"
)

func (r *episodeResolver) Titles(ctx context.Context, obj *models.Episode, first *int, skip *int) ([]*models.Title, error) {
	return sliceTitles(obj.Titles, first, skip), nil
}

func (r *episodeResolver) Synopses(ctx context.Context, obj *models.Episode, first *int, skip *int) ([]*models.Title, error) {
	return sliceTitles(obj.Synopses, first, skip), nil
}

func (r *episodeSetResolver) Media(ctx context.Context, obj *models.EpisodeSet) (*models.Media, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var md *models.Media
	err = ds.Database.Transaction(false, func(tx db.Tx) error {
		ser := ds.MediaService
		md, err = ser.GetByID(obj.MediaID, tx)
		if err != nil {
			return fmt.Errorf("failed to get Media by id %d: %w", obj.MediaID, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return md, nil
}

func (r *episodeSetResolver) Descriptions(ctx context.Context, obj *models.EpisodeSet, first *int, skip *int) ([]*models.Title, error) {
	return sliceTitles(obj.Descriptions, first, skip), nil
}

func (r *episodeSetResolver) Episodes(ctx context.Context, obj *models.EpisodeSet, first *int) ([]*models.Episode, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var list []*models.Episode
	err = ds.Database.Transaction(false, func(tx db.Tx) error {
		ser := ds.EpisodeService
		list, err = ser.GetMultiple(obj.Episodes, tx, nil)
		if err != nil {
			return fmt.Errorf("failed to get Epiosodes by ids: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}

// Episode returns EpisodeResolver implementation.
func (r *Resolver) Episode() EpisodeResolver { return &episodeResolver{r} }

// EpisodeSet returns EpisodeSetResolver implementation.
func (r *Resolver) EpisodeSet() EpisodeSetResolver { return &episodeSetResolver{r} }

type episodeResolver struct{ *Resolver }
type episodeSetResolver struct{ *Resolver }
