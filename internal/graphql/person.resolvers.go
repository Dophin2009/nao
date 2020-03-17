package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/Dophin2009/nao/pkg/db"
	"github.com/Dophin2009/nao/pkg/models"
)

func (r *personResolver) Names(ctx context.Context, obj *models.Person, first *int, skip *int) ([]*models.Title, error) {
	return sliceTitles(obj.Names, first, skip), nil
}

func (r *personResolver) Information(ctx context.Context, obj *models.Person, first *int, skip *int) ([]*models.Title, error) {
	return sliceTitles(obj.Information, first, skip), nil
}

func (r *personResolver) Media(ctx context.Context, obj *models.Person, first *int, skip *int) ([]*models.MediaCharacter, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var list []*models.MediaCharacter
	err = ds.Database.Transaction(false, func(tx db.Tx) error {
		ser := ds.MediaCharacterService
		list, err = ser.GetByPerson(obj.Meta.ID, first, skip, tx)
		if err != nil {
			return fmt.Errorf(
				"failed to get MediaCharacters by Person id %d: %w", obj.Meta.ID, err)
		}
		return nil
	})

	return list, nil
}

// Person returns PersonResolver implementation.
func (r *Resolver) Person() PersonResolver { return &personResolver{r} }

type personResolver struct{ *Resolver }
