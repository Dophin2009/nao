package graphql

import (
	"context"
	"fmt"

	"gitlab.com/Dophin2009/nao/pkg/data"
	"gitlab.com/Dophin2009/nao/pkg/db"
)

// TODO: Implement authentication

// Resolver is the root GraphQL resolver object.
type Resolver struct{}

func resolveMediaByID(ctx context.Context, mID int) (*data.Media, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var md *data.Media
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
	objTitles []data.Title, first *int, skip *int,
) []*data.Title {
	start, end := calculatePaginationBounds(first, skip, len(objTitles))

	titles := objTitles[start:end]
	tlist := make([]*data.Title, len(titles))
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

const (
	errmsgGetDataServices = "failed to get data services"
)

func errorGetDataServices(err error) error {
	return fmt.Errorf("failed to get data services: %w", err)
}
