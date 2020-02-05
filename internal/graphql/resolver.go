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

// Query returns a new QueryResolver.
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

// Mutation returns a new MutationResolver.
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// Character returns a new CharacterResolver.
func (r *Resolver) Character() CharacterResolver {
	return &characterResolver{r}
}

// Episode returns a new EpisodeResolver.
func (r *Resolver) Episode() EpisodeResolver {
	return &episodeResolver{r}
}

// EpisodeSet returns a new EpisodeSetResolver.
func (r *Resolver) EpisodeSet() EpisodeSetResolver {
	return &episodeSetResolver{r}
}

// Genre returns a new GenreResolver.
func (r *Resolver) Genre() GenreResolver {
	return &genreResolver{r}
}

// Media returns a new MediaResolver.
func (r *Resolver) Media() MediaResolver {
	return &mediaResolver{r}
}

// MediaCharacter returns a new MediaCharacterResolver.
func (r *Resolver) MediaCharacter() MediaCharacterResolver {
	return &mediaCharacterResolver{r}
}

// MediaGenre returns a new MediaGenreResolver.
func (r *Resolver) MediaGenre() MediaGenreResolver {
	return &mediaGenreResolver{r}
}

// MediaProducer returns a new MediaProducerResolver.
func (r *Resolver) MediaProducer() MediaProducerResolver {
	return &mediaProducerResolver{r}
}

// MediaRelation returns a new MediaRelationResolver.
func (r *Resolver) MediaRelation() MediaRelationResolver {
	return &mediaRelationResolver{r}
}

// Person returns a new PersonResolver.
func (r *Resolver) Person() PersonResolver {
	return &personResolver{r}
}

// Producer returns a new ProducerResolver.
func (r *Resolver) Producer() ProducerResolver {
	return &producerResolver{r}
}

// queryResolver is the root query resolver.
type queryResolver struct{ *Resolver }

// MediaByID resolves the query for the Media of the given ID.
func (r *queryResolver) MediaByID(
	ctx context.Context, id int,
) (*data.Media, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var md *data.Media
	err = ds.Database.Transaction(false, func(tx db.Tx) error {
		ser := ds.MediaService
		md, err = ser.GetByID(id, tx)
		if err != nil {
			return fmt.Errorf("failed to get Media by id %d: %w", id, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return md, nil
}

// mutationResolver is the root mutation resolver.
type mutationResolver struct{ *Resolver }

// CreateMedia resolves the mutation for creating a new Media.
func (r *mutationResolver) CreateMedia(
	ctx context.Context, media data.Media,
) (*data.Media, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	err = ds.Database.Transaction(true, func(tx db.Tx) error {
		ser := ds.MediaService
		_, err = ser.Create(&media, tx)
		if err != nil {
			return fmt.Errorf("failed to create Media: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &media, nil
}

// characterResolver is the field resolver for Character objects.
type characterResolver struct{ *Resolver }

// Names resolves the list of names for Character objects.
func (r *characterResolver) Names(
	ctx context.Context, obj *data.Character, first *int, skip *int,
) ([]*data.Title, error) {
	return sliceTitles(obj.Names, first, skip), nil
}

// Information resolves the list of information for Character objects.
func (r *characterResolver) Information(
	ctx context.Context, obj *data.Character, first *int, skip *int,
) ([]*data.Title, error) {
	return sliceTitles(obj.Information, first, skip), nil
}

// Media resolves the MediaCharacter list for Character objects.
func (r *characterResolver) Media(
	ctx context.Context, obj *data.Character, first *int, skip *int,
) ([]*data.MediaCharacter, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var list []*data.MediaCharacter
	err = ds.Database.Transaction(false, func(tx db.Tx) error {
		ser := ds.MediaCharacterService
		list, err = ser.GetByCharacter(obj.Meta.ID, first, skip, tx)
		if err != nil {
			return fmt.Errorf("failed to get MediaCharacters by Character id %d: %w",
				obj.Meta.ID, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}

// episodeResolver is the field resolver for Episode objects.
type episodeResolver struct{ *Resolver }

// Titles resolves the list of titles for Episode objects.
func (r *episodeResolver) Titles(
	ctx context.Context, obj *data.Episode, first *int, skip *int,
) ([]*data.Title, error) {
	return sliceTitles(obj.Titles, first, skip), nil
}

// Synopses resolves the list of synopses for Episode objects.
func (r *episodeResolver) Synopses(
	ctx context.Context, obj *data.Episode, first *int, skip *int,
) ([]*data.Title, error) {
	return sliceTitles(obj.Synopses, first, skip), nil
}

// episodeSetResolver is the field resolver for EpisodeSet objects.
type episodeSetResolver struct{ *Resolver }

// Descriptions resolves the list of descriptions for EpisodeSet objects.
func (r *episodeSetResolver) Descriptions(
	ctx context.Context, obj *data.EpisodeSet, first *int, skip *int,
) ([]*data.Title, error) {
	return sliceTitles(obj.Descriptions, first, skip), nil
}

// Media resolves the Media the EpisodeSet object belongs to.
func (r *episodeSetResolver) Media(
	ctx context.Context, obj *data.EpisodeSet) (*data.Media, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var md *data.Media
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

// Episodes resolves the Episode list for EpisodeSet objects.
func (r *episodeSetResolver) Episodes(
	ctx context.Context, obj *data.EpisodeSet, first *int, skip *int,
) ([]*data.Episode, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var list []*data.Episode
	err = ds.Database.Transaction(false, func(tx db.Tx) error {
		ser := ds.EpisodeService
		list, err = ser.GetMultiple(obj.Episodes, first, skip, tx, nil)
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

// genreResolver is the field resolver for Genre objects.
type genreResolver struct{ *Resolver }

// Names resolves the list of names for Genre objects.
func (r *genreResolver) Names(
	ctx context.Context, obj *data.Genre, first *int, skip *int,
) ([]*data.Title, error) {
	return sliceTitles(obj.Names, first, skip), nil
}

// Descriptions resolves the list of descriptions for Genre objects.
func (r *genreResolver) Descriptions(
	ctx context.Context, obj *data.Genre, first *int, skip *int,
) ([]*data.Title, error) {
	return sliceTitles(obj.Descriptions, first, skip), nil
}

// Media resolves the Media in the relationship for MediaGenre objects.
func (r *genreResolver) Media(
	ctx context.Context, obj *data.Genre, first *int, skip *int,
) ([]*data.MediaGenre, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var list []*data.MediaGenre
	err = ds.Database.Transaction(false, func(tx db.Tx) error {
		ser := ds.MediaGenreService
		list, err = ser.GetByGenre(obj.Meta.ID, first, skip, tx)
		if err != nil {
			return fmt.Errorf("failed to get MediaGenres by Genre id %d: %w",
				obj.Meta.ID, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}

// mediaResolver is the field resolver for Media objects.
type mediaResolver struct{ *Resolver }

// Titles resolves the title list for Media objects.
func (r *mediaResolver) Titles(
	ctx context.Context, obj *data.Media, first *int, skip *int,
) ([]*data.Title, error) {
	return sliceTitles(obj.Titles, first, skip), nil
}

// Synopses resolves the synopses list for Media objects.
func (r *mediaResolver) Synopses(
	ctx context.Context, obj *data.Media, first *int, skip *int,
) ([]*data.Title, error) {
	return sliceTitles(obj.Synopses, first, skip), nil
}

// Background resolves the background information lists for Media objects.
func (r *mediaResolver) Background(
	ctx context.Context, obj *data.Media, first *int, skip *int,
) ([]*data.Title, error) {
	return sliceTitles(obj.Titles, first, skip), nil
}

// EpisodeSets resolves the EpisodeSets for Media objects.
func (r *mediaResolver) EpisodeSets(
	ctx context.Context, obj *data.Media, first *int, skip *int,
) ([]*data.EpisodeSet, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var list []*data.EpisodeSet
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

// Producers resolves the MediaProducer relationships for Media objects.
func (r *mediaResolver) Producers(
	ctx context.Context, obj *data.Media, first *int, skip *int,
) ([]*data.MediaProducer, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var list []*data.MediaProducer
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

// Characters resolves the MediaCharacter relationships for Media objects.
func (r *mediaResolver) Characters(
	ctx context.Context, obj *data.Media, first *int, skip *int,
) ([]*data.MediaCharacter, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var list []*data.MediaCharacter
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

// Genres resolves the MediaGenre relationships for Media objects.
func (r *mediaResolver) Genres(
	ctx context.Context, obj *data.Media, first *int, skip *int,
) ([]*data.MediaGenre, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var list []*data.MediaGenre
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

// mediaCharacterResolver is the field resolver for MediaCharacter objects.
type mediaCharacterResolver struct{ *Resolver }

// Media resolves the Media in the relationship for MediaCharacter objects.
func (r *mediaCharacterResolver) Media(
	ctx context.Context, obj *data.MediaCharacter,
) (*data.Media, error) {
	return resolveMediaByID(ctx, obj.MediaID)
}

// Character resolves the Character in the relationship for MediaCharacter
// objects.
func (r *mediaCharacterResolver) Character(
	ctx context.Context, obj *data.MediaCharacter,
) (*data.Character, error) {
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

// Person resolves the Person in the relationship for MediaCharacter objects.
func (r *mediaCharacterResolver) Person(
	ctx context.Context, obj *data.MediaCharacter,
) (*data.Person, error) {
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

// mediaGenreResolver is the field resolver for MediaGenre objects.
type mediaGenreResolver struct{ *Resolver }

// Media resolves the Media in the relationship for MediaGenre objects.
func (r *mediaGenreResolver) Media(
	ctx context.Context, obj *data.MediaGenre,
) (*data.Media, error) {
	return resolveMediaByID(ctx, obj.MediaID)
}

// Genre resolves the Genre in the relationship for MediaGenre objects.
func (r *mediaGenreResolver) Genre(
	ctx context.Context, obj *data.MediaGenre,
) (*data.Genre, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var g *data.Genre
	err = ds.Database.Transaction(false, func(tx db.Tx) error {
		ser := ds.GenreService
		g, err = ser.GetByID(obj.GenreID, tx)
		if err != nil {
			return fmt.Errorf("failed to get Genre by id %d: %w", obj.GenreID, err)
		}
		return nil
	})

	return g, nil
}

// mediaProducerResolver is the field resolver for MediaProducer objects.
type mediaProducerResolver struct{ *Resolver }

// Media resolves the Media in the relationship for MediaProducer objects.
func (r *mediaProducerResolver) Media(
	ctx context.Context, obj *data.MediaProducer,
) (*data.Media, error) {
	return resolveMediaByID(ctx, obj.MediaID)
}

// Producer resolves the Producer in the relationship for MediaProducer
// objects.
func (r *mediaProducerResolver) Producer(
	ctx context.Context, obj *data.MediaProducer,
) (*data.Producer, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var p *data.Producer
	err = ds.Database.Transaction(false, func(tx db.Tx) error {
		ser := ds.ProducerService
		p, err = ser.GetByID(obj.ProducerID, tx)
		if err != nil {
			return fmt.Errorf("failed to get Producer by id %d: %w", obj.ProducerID, err)
		}
		return nil
	})

	return p, nil
}

// mediaRelationResolver is the field resolver for MediaRelation objects.
type mediaRelationResolver struct{ *Resolver }

// Owner resolves the owning Media in the relationship for MediaRelation
// objects.
func (r *mediaRelationResolver) Owner(
	ctx context.Context, obj *data.MediaRelation,
) (*data.Media, error) {
	return resolveMediaByID(ctx, obj.OwnerID)
}

// Related resolves the owned Media in the relationship for MediaRelation
// objects.
func (r *mediaRelationResolver) Related(
	ctx context.Context, obj *data.MediaRelation,
) (*data.Media, error) {
	return resolveMediaByID(ctx, obj.RelatedID)
}

// personResolver is the field resolver for Person objects.
type personResolver struct{ *Resolver }

// Names resolves the list of names for Person objects.
func (r *personResolver) Names(
	ctx context.Context, obj *data.Person, first *int, skip *int,
) ([]*data.Title, error) {
	return sliceTitles(obj.Names, first, skip), nil
}

// Information resolves the list of information segments for Person objects.
func (r *personResolver) Information(
	ctx context.Context, obj *data.Person, first *int, skip *int,
) ([]*data.Title, error) {
	return sliceTitles(obj.Information, first, skip), nil
}

// Media resolves the MediaCharacter relationships for
// Person objects.
func (r *personResolver) Media(
	ctx context.Context, obj *data.Person, first *int, skip *int,
) ([]*data.MediaCharacter, error) {
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

// producerResolver is the field resolver for Producer objects.
type producerResolver struct{ *Resolver }

// Titles resolves the list of titles for Producer objects.
func (r *producerResolver) Titles(
	ctx context.Context, obj *data.Producer, first *int, skip *int,
) ([]*data.Title, error) {
	return sliceTitles(obj.Titles, first, skip), nil
}

// Media resolves the Media list for Producer objects.
func (r *producerResolver) Media(
	ctx context.Context, obj *data.Producer, first *int, skip *int,
) ([]*data.MediaProducer, error) {
	ds, err := getCtxDataService(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	var list []*data.MediaProducer
	err = ds.Database.Transaction(false, func(tx db.Tx) error {
		ser := ds.MediaProducerService
		list, err = ser.GetByProducer(obj.Meta.ID, first, skip, tx)
		if err != nil {
			return fmt.Errorf(
				"failed to get MediaProducers by Producer id %d: %w", obj.Meta.ID, err)
		}
		return nil
	})

	return list, nil
}

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
