package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.com/Dophin2009/nao/pkg/data"
	"gitlab.com/Dophin2009/nao/pkg/db"
)

func (r *personResolver) Names(ctx context.Context, obj *data.Person, first *int, skip *int) ([]*data.Title, error) {
	return sliceTitles(obj.Names, first, skip), nil
}

func (r *personResolver) Information(ctx context.Context, obj *data.Person, first *int, skip *int) ([]*data.Title, error) {
	return sliceTitles(obj.Information, first, skip), nil
}

func (r *personResolver) Media(ctx context.Context, obj *data.Person, first *int, skip *int) ([]*data.MediaCharacter, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var list []*data.MediaCharacter
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
