package naos

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gitlab.com/Dophin2009/nao/internal/data"
	"gitlab.com/Dophin2009/nao/internal/graphql"
	"gitlab.com/Dophin2009/nao/internal/web"
	bolt "go.etcd.io/bbolt"
)

// Application is the main naos application.
type Application struct {
	Server *web.Server
	DB     *bolt.DB
}

// HTTPServer returns the application's HTTP server.
func (a *Application) HTTPServer() http.Server {
	return a.Server.HTTPServer()
}

// NewApplication returns a new naos Application.
func NewApplication(c *Configuration) (*Application, error) {
	// Open database connection
	log.Println("Establishing database connection")
	db, err := data.ConnectDatabase(
		c.DB.Path, os.FileMode(c.DB.Filemode), true,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create the API controller and HTTP server
	address := fmt.Sprintf("%s:%s", c.Hostname, c.Port)
	s := web.NewServer(address)

	ds := graphql.DataServices{
		CharacterService:      &data.CharacterService{DB: db},
		EpisodeService:        &data.EpisodeService{DB: db},
		EpisodeSetService:     &data.EpisodeSetService{DB: db},
		GenreService:          &data.GenreService{DB: db},
		MediaService:          &data.MediaService{DB: db},
		MediaCharacterService: &data.MediaCharacterService{DB: db},
		MediaGenreService:     &data.MediaGenreService{DB: db},
		MediaProducerService:  &data.MediaProducerService{DB: db},
		MediaRelationSerivce:  &data.MediaRelationService{DB: db},
		PersonService:         &data.PersonService{DB: db},
		ProducerService:       &data.ProducerService{DB: db},
		UserService:           &data.UserService{DB: db},
		UserMediaService:      &data.UserMediaService{DB: db},
		UserMediaListService:  &data.UserMediaListService{DB: db},
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
		Server: &s,
		DB:     db,
	}, nil
}
