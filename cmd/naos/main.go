package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"gitlab.com/Dophin2009/nao/internal/data"
	"gitlab.com/Dophin2009/nao/internal/naos"
	"gitlab.com/Dophin2009/nao/internal/naos/graphql"
	"gitlab.com/Dophin2009/nao/internal/web"
	bolt "go.etcd.io/bbolt"
)

func main() {
	// Exit with status code 0 at the end
	defer os.Exit(0)

	println("-------------------: NAO SERVER :-------------------")

	// Read configuration files
	conf, err := naos.ReadConfigs()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	// Open database connection
	log.Println("Establishing database connection")
	db, err := data.ConnectDatabase(conf.DB.Path, os.FileMode(conf.DB.Filemode), true)
	if err != nil {
		log.Fatal("Error connecting to database ")
		return
	}
	// Clear database and close connection at the end
	defer db.Close()
	defer data.ClearDatabase(db)

	// Create the API controller and HTTP server
	address := fmt.Sprintf("%s:%s", conf.Hostname, conf.Port)
	s, err := initServer(address, db)
	if err != nil {
		log.Fatalf("Error initializing server: %v", err)
		return
	}
	shttp := s.HTTPServer()

	// Launch server in goroutine
	go func() {
		log.Println("Launching server on", s.Address)
		err := shttp.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for SIGINTERRUPT signal
	wait := time.Second * 15
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc

	// Wait for processes to end, then shutdown
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	shttp.Shutdown(ctx)

	println()
	log.Println("Exiting...")
}

func initServer(address string, db *bolt.DB) (*web.Server, error) {
	s := web.NewServer(address)

	ds := graphql.DataServices{
		CharacterService:      &data.CharacterService{DB: db},
		EpisodeService:        &data.EpisodeService{DB: db},
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

	graphqlHandler := naos.NewGraphQLHandler([]string{"graphql"}, &ds)
	s.RegisterHandler(graphqlHandler)

	graphiqlHandler, err := naos.NewGraphiQLHandler(
		[]string{"graphiql"}, graphqlHandler.PathString(),
	)
	if err != nil {
		return nil, err
	}
	s.RegisterHandler(graphiqlHandler)

	playgroundHandler := naos.NewGraphQLPlaygroundHandler([]string{"playground"}, graphqlHandler.PathString())
	s.RegisterHandler(playgroundHandler)

	return &s, nil
}
