package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/Dophin2009/nao/pkg/data"
	"github.com/Dophin2009/nao/pkg/db"
)

func (r *mediaCharacterResolver) Media(ctx context.Context, obj *data.MediaCharacter) (*data.Media, error) {
	return resolveMediaByID(ctx, obj.MediaID)
}

func (r *mediaCharacterResolver) Character(ctx context.Context, obj *data.MediaCharacter) (*data.Character, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	if obj.CharacterID == nil {
		return nil, nil
	}

	var c *data.Character
	err = ds.Database.Transaction(false, func(tx db.Tx) error {
		ser := ds.CharacterService
		c, err = ser.GetByID(*obj.CharacterID, tx)
		if err != nil {
			return fmt.Errorf(
				"failed to get Character by id %d: %w", *obj.CharacterID, err)
		}
		return nil
	})

	return c, nil
}

func (r *mediaCharacterResolver) Person(ctx context.Context, obj *data.MediaCharacter) (*data.Person, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	if obj.PersonID == nil {
		return nil, nil
	}

	var p *data.Person
	err = ds.Database.Transaction(false, func(tx db.Tx) error {
		ser := ds.PersonService
		p, err = ser.GetByID(*obj.PersonID, tx)
		if err != nil {
			return fmt.Errorf(
				"failed to get Person by id %d: %w", *obj.PersonID, err)
		}
		return nil
	})

	return p, nil
}

// MediaCharacter returns MediaCharacterResolver implementation.
func (r *Resolver) MediaCharacter() MediaCharacterResolver { return &mediaCharacterResolver{r} }

type mediaCharacterResolver struct{ *Resolver }
