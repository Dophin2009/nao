package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/Dophin2009/nao/pkg/db"
	"github.com/Dophin2009/nao/pkg/models"
)

func (r *mediaResolver) Titles(ctx context.Context, obj *models.Media, first *int, skip *int) ([]*models.Title, error) {
	return sliceTitles(obj.Titles, first, skip), nil
}

func (r *mediaResolver) Synopses(ctx context.Context, obj *models.Media, first *int, skip *int) ([]*models.Title, error) {
	return sliceTitles(obj.Synopses, first, skip), nil
}

func (r *mediaResolver) Background(ctx context.Context, obj *models.Media, first *int, skip *int) ([]*models.Title, error) {
	return sliceTitles(obj.Titles, first, skip), nil
}

func (r *mediaResolver) EpisodeSets(ctx context.Context, obj *models.Media, first *int, skip *int) ([]*models.EpisodeSet, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var list []*models.EpisodeSet
	ds.Database.Transaction(false, func(tx db.Tx) error {
		ser := ds.EpisodeSetService
		list, err = ser.GetByMedia(obj.Meta.ID, first, skip, tx)
		if err != nil {
			return fmt.Errorf("failed to get EpisodeSets by Media id %d: %w",
				obj.Meta.ID, err)
		}
		return nil
	})

	return list, nil
}

func (r *mediaResolver) Producers(ctx context.Context, obj *models.Media, first *int, skip *int) ([]*models.MediaProducer, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var list []*models.MediaProducer
	err = ds.Database.Transaction(false, func(tx db.Tx) error {
		ser := ds.MediaProducerService
		list, err = ser.GetByMedia(obj.Meta.ID, first, skip, tx)
		if err != nil {
			return fmt.Errorf("failed to get MediaProducers by Media id %d: %w",
				obj.Meta.ID, err)
		}
		return nil
	})

	return list, nil
}

func (r *mediaResolver) Characters(ctx context.Context, obj *models.Media, first *int, skip *int) ([]*models.MediaCharacter, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var list []*models.MediaCharacter
	err = ds.Database.Transaction(false, func(tx db.Tx) error {
		ser := ds.MediaCharacterService
		list, err = ser.GetByMedia(obj.Meta.ID, first, skip, tx)
		if err != nil {
			return fmt.Errorf(
				"failed to get MediaCharacters by Media id %d: %w", obj.Meta.ID, err)
		}
		return nil
	})

	return list, nil
}

func (r *mediaResolver) Genres(ctx context.Context, obj *models.Media, first *int, skip *int) ([]*models.MediaGenre, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var list []*models.MediaGenre
	err = ds.Database.Transaction(false, func(tx db.Tx) error {
		ser := ds.MediaGenreService
		list, err = ser.GetByMedia(obj.Meta.ID, first, skip, tx)
		if err != nil {
			return fmt.Errorf("failed to get MediaGenres by Media id %d: %w",
				obj.Meta.ID, err)
		}
		return nil
	})

	return list, nil
}

// Media returns MediaResolver implementation.
func (r *Resolver) Media() MediaResolver { return &mediaResolver{r} }

type mediaResolver struct{ *Resolver }
