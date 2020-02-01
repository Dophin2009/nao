package graphql

import (
	"context"
	"errors"

	"gitlab.com/Dophin2009/nao/internal/data"
)

// DataServices contains all data layer services required,
// to be passed around in a context object.
type DataServices struct {
	CharacterService      *data.CharacterService
	EpisodeService        *data.EpisodeService
	EpisodeSetService     *data.EpisodeSetService
	GenreService          *data.GenreService
	MediaService          *data.MediaService
	MediaCharacterService *data.MediaCharacterService
	MediaGenreService     *data.MediaGenreService
	MediaProducerService  *data.MediaProducerService
	MediaRelationSerivce  *data.MediaRelationService
	PersonService         *data.PersonService
	ProducerService       *data.ProducerService
	UserService           *data.UserService
	UserMediaService      *data.UserMediaService
	UserMediaListService  *data.UserMediaListService
}

// DataServicesKey is the context key value for DataServices.
const DataServicesKey = "DataServicesKey"

func getDataServicesFromCtx(ctx context.Context) (*DataServices, error) {
	v, ok := ctx.Value(DataServicesKey).(*DataServices)
	if !ok {
		return nil, errors.New("DataServices not found in context")
	}
	return v, nil
}
