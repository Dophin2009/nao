package graphql

import (
	"context"
	"errors"

	"gitlab.com/Dophin2009/nao/pkg/data"
	"gitlab.com/Dophin2009/nao/pkg/db"
)

// DataService contains all data layer services required, to be passed around
// in a context object.
type DataService struct {
	Database              db.DatabaseService
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

// DataServiceKey is the context key value for DataServices.
const DataServiceKey = "DataServicesKey"

func getCtxDataService(ctx context.Context) (*DataService, error) {
	v, ok := ctx.Value(DataServiceKey).(*DataService)
	if !ok {
		return nil, errors.New("DataServices not found in context")
	}
	return v, nil
}
