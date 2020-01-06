package graphql

import (
	"context"
	"fmt"

	"gitlab.com/Dophin2009/nao/internal/data"
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

// MediaByID resolves the query for the Media of the given
// ID.
func (r *queryResolver) MediaByID(ctx context.Context, id int) (*data.Media, error) {
	ds, err := getDataServicesFromCtx(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	ser := ds.MediaService
	md, err := ser.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get Media by id %d: %w", id, err)
	}

	return md, nil
}

// mutationResolver is the root mutation resolver.
type mutationResolver struct{ *Resolver }

// CreateMedia resolves the mutation for creating a new Media.
func (r *mutationResolver) CreateMedia(ctx context.Context, media data.Media) (*data.Media, error) {
	ds, err := getDataServicesFromCtx(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	ser := ds.MediaService
	err = ser.Create(&media)
	if err != nil {
		return nil, fmt.Errorf("failed to create Media: %w", err)
	}

	return &media, nil
}

// characterResolver is the field resolver for Character objects.
type characterResolver struct{ *Resolver }

// Media resolves the MediaCharacter list for Character objects.
func (r *characterResolver) Media(ctx context.Context, obj *data.Character, first int, prefix *int) ([]*data.MediaCharacter, error) {
	ds, err := getDataServicesFromCtx(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	ser := ds.MediaCharacterService
	list, err := ser.GetByCharacter(obj.ID, first, prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to get MediaCharacters by Character id %d: %w", obj.ID, err)
	}

	return list, nil
}

// episodeResolver is the field resolver for Episode objects.
type episodeResolver struct{ *Resolver }

// Media resolves the parent Media for Episode objects.
func (r *episodeResolver) Media(ctx context.Context, obj *data.Episode) (*data.Media, error) {
	return resolveMediaByID(ctx, obj.MediaID)
}

// genreResolver is the field resolver for Genre objects.
type genreResolver struct{ *Resolver }

// Media resolves the Media in the relationship for
// MediaGenre objects.
func (r *genreResolver) Media(ctx context.Context, obj *data.Genre, first int, prefix *int) ([]*data.MediaGenre, error) {
	ds, err := getDataServicesFromCtx(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	ser := ds.MediaGenreService
	list, err := ser.GetByGenre(obj.ID, first, prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to get MediaGenres by Genre id %d: %w", obj.ID, err)
	}

	return list, nil
}

// mediaResolver is the field resolver for Media objects.
type mediaResolver struct{ *Resolver }

// Episodes resolves the Episodes for Media objects.
func (r *mediaResolver) Episodes(ctx context.Context, obj *data.Media, first int, prefix *int) ([]*data.Episode, error) {
	ds, err := getDataServicesFromCtx(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	ser := ds.EpisodeService
	list, err := ser.GetByMedia(obj.ID, first, prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to get Episodes by Media id %d: %w", obj.ID, err)
	}

	return list, nil
}

// Producers resolves the MediaProducer relationships for
// Media objects.
func (r *mediaResolver) Producers(ctx context.Context, obj *data.Media, first int, prefix *int) ([]*data.MediaProducer, error) {
	ds, err := getDataServicesFromCtx(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	ser := ds.MediaProducerService
	list, err := ser.GetByMedia(obj.ID, first, prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to get MediaProducers by Media id %d: %w", obj.ID, err)
	}

	return list, nil
}

// Characters resolves the MediaCharacter relationships for
// Media objects.
func (r *mediaResolver) Characters(ctx context.Context, obj *data.Media, first int, prefix *int) ([]*data.MediaCharacter, error) {
	ds, err := getDataServicesFromCtx(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	ser := ds.MediaCharacterService
	list, err := ser.GetByMedia(obj.ID, first, prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to get MediaCharacters by Media id %d: %w", obj.ID, err)
	}

	return list, nil
}

// Genres resolves the MediaGenre relationships for
// Media objects.
func (r *mediaResolver) Genres(ctx context.Context, obj *data.Media, first int, prefix *int) ([]*data.MediaGenre, error) {
	ds, err := getDataServicesFromCtx(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	ser := ds.MediaGenreService
	list, err := ser.GetByMedia(obj.ID, first, prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to get MediaGenres by Media id %d: %w", obj.ID, err)
	}

	return list, nil
}

// mediaCharacterResolver is the field resolver for
// MediaCharacter objects.
type mediaCharacterResolver struct{ *Resolver }

// Media resolves the Media in the relationship for
// MediaCharacter objects.
func (r *mediaCharacterResolver) Media(ctx context.Context, obj *data.MediaCharacter) (*data.Media, error) {
	return resolveMediaByID(ctx, obj.MediaID)
}

// Character resolves the Character in the relationship for
// MediaCharacter objects.
func (r *mediaCharacterResolver) Character(ctx context.Context, obj *data.MediaCharacter) (*data.Character, error) {
	ds, err := getDataServicesFromCtx(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	ser := ds.CharacterService
	if obj.CharacterID == nil {
		return nil, nil
	}

	c, err := ser.GetByID(*obj.CharacterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get Character by id %d: %w", *obj.CharacterID, err)
	}

	return c, nil
}

// Person resolves the Person in the relationship for
// MediaCharacter objects
func (r *mediaCharacterResolver) Person(ctx context.Context, obj *data.MediaCharacter) (*data.Person, error) {
	ds, err := getDataServicesFromCtx(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	ser := ds.PersonService
	if obj.PersonID == nil {
		return nil, nil
	}

	p, err := ser.GetByID(*obj.PersonID)
	if err != nil {
		return nil, fmt.Errorf("failed to get Person by id %d: %w", *obj.PersonID, err)
	}

	return p, nil
}

// mediaGenreResolver is the field resolver for MediaGenre
// objects.
type mediaGenreResolver struct{ *Resolver }

// Media resolves the Media in the relationship for MediaGenre
// objects.
func (r *mediaGenreResolver) Media(ctx context.Context, obj *data.MediaGenre) (*data.Media, error) {
	return resolveMediaByID(ctx, obj.MediaID)
}

// Genre resolves the Genre in the relationship for MediaGenre
// objects.
func (r *mediaGenreResolver) Genre(ctx context.Context, obj *data.MediaGenre) (*data.Genre, error) {
	ds, err := getDataServicesFromCtx(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	ser := ds.GenreService
	g, err := ser.GetByID(obj.GenreID)
	if err != nil {
		return nil, fmt.Errorf("failed to get Genre by id %d: %w", obj.GenreID, err)
	}

	return g, nil
}

// mediaProducerResolver is the field resolver for MediaProducer
// objects.
type mediaProducerResolver struct{ *Resolver }

// Media resolves the Media in the relationship for
// MediaProducer objects.
func (r *mediaProducerResolver) Media(ctx context.Context, obj *data.MediaProducer) (*data.Media, error) {
	return resolveMediaByID(ctx, obj.MediaID)
}

// Producer resolves the Producer in the relationship for
// MediaProducer objects.
func (r *mediaProducerResolver) Producer(ctx context.Context, obj *data.MediaProducer) (*data.Producer, error) {
	ds, err := getDataServicesFromCtx(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	ser := ds.ProducerService
	p, err := ser.GetByID(obj.ProducerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get Producer by id %d: %w", obj.ProducerID, err)
	}

	return p, nil
}

// mediaRelationResolver is the field resolver for
// MediaRelation objects.
type mediaRelationResolver struct{ *Resolver }

// Owner resolves the owning Media in the relationship for
// MediaRelation objects.
func (r *mediaRelationResolver) Owner(ctx context.Context, obj *data.MediaRelation) (*data.Media, error) {
	return resolveMediaByID(ctx, obj.OwnerID)
}

// Related resolves the owned Media in the relationship for
// MediaRelation objects.
func (r *mediaRelationResolver) Related(ctx context.Context, obj *data.MediaRelation) (*data.Media, error) {
	return resolveMediaByID(ctx, obj.RelatedID)
}

// personResolver is the field resolver for Person
// objects.
type personResolver struct{ *Resolver }

// Media resolves the MediaCharacter relationships for
// Person objects.
func (r *personResolver) Media(ctx context.Context, obj *data.Person, first int, prefix *int) ([]*data.MediaCharacter, error) {
	ds, err := getDataServicesFromCtx(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	ser := ds.MediaCharacterService
	list, err := ser.GetByPerson(obj.ID, first, prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to get MediaCharacters by Person id %d: %w", obj.ID, err)
	}

	return list, nil
}

// producerResolver is the field resolver for Producer
// objects.
type producerResolver struct{ *Resolver }

func (r *producerResolver) Media(ctx context.Context, obj *data.Producer, first int, prefix *int) ([]*data.MediaProducer, error) {
	ds, err := getDataServicesFromCtx(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	ser := ds.MediaProducerService
	list, err := ser.GetByProducer(obj.ID, first, prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to get MediaProducers by Producer id %d: %w", obj.ID, err)
	}

	return list, nil
}

func resolveMediaByID(ctx context.Context, mID int) (*data.Media, error) {
	ds, err := getDataServicesFromCtx(ctx)
	if err != nil {
		return nil, errorGetDataServices(err)
	}

	ser := ds.MediaService
	md, err := ser.GetByID(mID)
	if err != nil {
		return nil, fmt.Errorf("failed to get Media by id %d: %w", mID, err)
	}

	return md, nil
}

const (
	errmsgGetDataServices = "failed to get data services"
)

func errorGetDataServices(err error) error {
	return fmt.Errorf("failed to get data services: %w", err)
}
