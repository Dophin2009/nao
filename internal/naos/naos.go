package naos

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"gitlab.com/Dophin2009/nao/internal/graphql"
	"gitlab.com/Dophin2009/nao/pkg/data"
	"gitlab.com/Dophin2009/nao/pkg/web"
)

// Application is the main naos application.
type Application struct {
	Server    *web.Server
	DataLayer *graphql.DataService
}

// HTTPServer returns the application's HTTP server.
func (a *Application) HTTPServer() http.Server {
	return a.Server.HTTPServer()
}

// NewApplication returns a new naos Application.
func NewApplication(c *Configuration) (*Application, error) {
	// Open database connection
	log.WithFields(log.Fields{
		"path":     c.DB.Path,
		"filemode": c.DB.Filemode,
	}).Info("Establishing database connection")

	// Create the API controller and HTTP server
	address := fmt.Sprintf("%s:%s", c.Hostname, c.Port)
	s := web.NewServer(address)

	characterService := &data.CharacterService{}
	episodeService := &data.EpisodeService{}
	episodeSetService := &data.EpisodeSetService{}
	genreService := &data.GenreService{}
	mediaService := &data.MediaService{}
	personService := &data.PersonService{}
	producerService := &data.ProducerService{}
	userService := &data.UserService{}

	mediaCharacterService := &data.MediaCharacterService{
		MediaService:     mediaService,
		CharacterService: characterService,
		PersonService:    personService,
	}
	mediaGenreService := &data.MediaGenreService{
		MediaService: mediaService,
		GenreService: genreService,
	}
	mediaProducerService := &data.MediaProducerService{
		MediaService:    mediaService,
		ProducerService: producerService,
	}
	mediaRelationService := &data.MediaRelationService{
		MediaService: mediaService,
	}
	userMediaService := &data.UserMediaService{
		UserService:  userService,
		MediaService: mediaService,
	}
	userMediaListService := &data.UserMediaListService{
		UserService:      userService,
		UserMediaService: userMediaService,
	}

	buckets := []string{
		characterService.Bucket(), episodeService.Bucket(), episodeSetService.Bucket(),
		genreService.Bucket(), mediaService.Bucket(), personService.Bucket(),
		producerService.Bucket(), userService.Bucket(), mediaCharacterService.Bucket(),
		mediaGenreService.Bucket(), mediaProducerService.Bucket(),
		mediaRelationService.Bucket(), userMediaService.Bucket(),
		userMediaListService.Bucket(),
	}

	db, err := data.ConnectBoltDatabase(&data.BoltDatabaseConfig{
		Path:         c.DB.Path,
		FileMode:     os.FileMode(c.DB.Filemode),
		Buckets:      buckets,
		ClearOnClose: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	ds := graphql.DataService{
		Database:              db,
		CharacterService:      characterService,
		EpisodeService:        episodeService,
		EpisodeSetService:     episodeSetService,
		GenreService:          genreService,
		MediaService:          mediaService,
		MediaCharacterService: mediaCharacterService,
		MediaGenreService:     mediaGenreService,
		MediaProducerService:  mediaProducerService,
		MediaRelationSerivce:  mediaRelationService,
		PersonService:         personService,
		ProducerService:       producerService,
		UserService:           userService,
		UserMediaService:      userMediaService,
		UserMediaListService:  userMediaListService,
	}

	graphqlHandler := NewGraphQLHandler([]string{"graphql"}, &ds)
	s.RegisterHandler(graphqlHandler)

	graphiqlHandler, err := NewGraphiQLHandler(
		[]string{"graphiql"}, graphqlHandler.PathString(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create GraphiQL handler: %w", err)
	}

	s.RegisterHandler(graphiqlHandler)

	return &Application{
		Server:    &s,
		DataLayer: &ds,
	}, nil
}
