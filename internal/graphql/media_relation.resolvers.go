package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.com/Dophin2009/nao/pkg/data"
)

func (r *mediaRelationResolver) Owner(ctx context.Context, obj *data.MediaRelation) (*data.Media, error) {
	return resolveMediaByID(ctx, obj.OwnerID)
}

func (r *mediaRelationResolver) Related(ctx context.Context, obj *data.MediaRelation) (*data.Media, error) {
	return resolveMediaByID(ctx, obj.RelatedID)
}

// MediaRelation returns MediaRelationResolver implementation.
func (r *Resolver) MediaRelation() MediaRelationResolver { return &mediaRelationResolver{r} }

type mediaRelationResolver struct{ *Resolver }
