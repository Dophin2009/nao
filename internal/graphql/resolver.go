package graphql

import (
	"context"
	"errors"
	"fmt"

	"github.com/Dophin2009/nao/pkg/data"
	"github.com/Dophin2009/nao/pkg/data/models"
	"github.com/Dophin2009/nao/pkg/db"
)

// TODO: Implement authentication

// Resolver is the root GraphQL resolver object.
type Resolver struct{}

func resolveMediaByID(ctx context.Context, mID int) (*models.Media, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var md *models.Media
	err = ds.Database.Transaction(false, func(tx db.Tx) error {
		ser := ds.MediaService
		md, err = ser.GetByID(mID, tx)
		if err != nil {
			return fmt.Errorf("failed to get Media by id %d: %w", mID, err)
		}
		return nil
	})

	return md, nil
}

func sliceTitles(
	objTitles []models.Title, first *int, skip *int,
) []*models.Title {
	start, end := calculatePaginationBounds(first, skip, len(objTitles))

	titles := objTitles[start:end]
	tlist := make([]*models.Title, len(titles))
	for i := range tlist {
		tlist[i] = &titles[i]
	}
	return tlist
}

func calculatePaginationBounds(first *int, skip *int, size int) (int, int) {
	if size <= 0 {
		return 0, 0
	}

	var start int
	if skip == nil || *skip <= 0 {
		start = 0
	} else {
		start = *skip
	}

	if start >= size {
		start = size
	}

	var end int
	if first == nil || *first < 0 {
		end = size
	} else {
		end = start + *first
	}

	if end > size {
		end = size
	}

	return start, end
}

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

const (
	errmsgGetDataServices = "failed to get data services"
)

func errorGetDataServices(err error) error {
	return fmt.Errorf("failed to get data services: %w", err)
}
