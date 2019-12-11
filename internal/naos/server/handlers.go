package server

import (
	"gitlab.com/Dophin2009/nao/pkg/data"
	bolt "go.etcd.io/bbolt"
)

// NewEntityHandlerGroups returns a list of all
// entity-related handler groups with
func NewEntityHandlerGroups(db *bolt.DB) []HandlerGroup {
	return []HandlerGroup{
		&MediaHandlerGroup{
			Service: &data.MediaService{DB: db},
		},
		&EpisodeHandlerGroup{
			Service: &data.EpisodeService{DB: db},
		},
		&CharacterHandlerGroup{
			Service: &data.CharacterService{DB: db},
		},
		&GenreHandlerGroup{
			Service: &data.GenreService{DB: db},
		},
		&ProducerHandlerGroup{
			Service: &data.ProducerService{DB: db},
		},
		&UserHandlerGroup{
			Service: &data.UserService{DB: db},
		},
		&MediaRelationHandlerGroup{
			Service: &data.MediaRelationService{DB: db},
		},
		&MediaCharacterHandlerGroup{
			Service: &data.MediaCharacterService{DB: db},
		},
		&MediaGenreHandlerGroup{
			Service: &data.MediaGenreService{DB: db},
		},
		&MediaProducerHandlerGroup{
			Service: &data.MediaProducerService{DB: db},
		},
		&UserMediaHandlerGroup{
			Service: &data.UserMediaService{DB: db},
		},
		&UserMediaListHandlerGroup{
			Service: &data.UserMediaListService{DB: db},
		},
	}
}
