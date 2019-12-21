package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"gitlab.com/Dophin2009/nao/internal/config"
	"gitlab.com/Dophin2009/nao/internal/data"
	"gitlab.com/Dophin2009/nao/internal/naos/gqlschema"
	"gitlab.com/Dophin2009/nao/internal/naos/handlers"
	"gitlab.com/Dophin2009/nao/internal/web"
	bolt "go.etcd.io/bbolt"
)

func main() {
	// Exit with status code 0 at the end
	defer os.Exit(0)

	println("-------------------: NAO SERVER :-------------------")

	// Read configuration files
	conf, err := config.ReadLinuxConfigs("nao")
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

	ds := &gqlschema.DataServices{
		MediaService: &data.MediaService{
			DB: db,
		},
	}

	schema, err := gqlschema.Schema()
	if err != nil {
		return nil, err
	}

	ctx := context.WithValue(context.Background(), gqlschema.ContextDataServices, ds)
	graphqlHandler := handlers.NewGraphQLHandler(ctx, &schema, []string{"graphql"})
	s.RegisterHandler(graphqlHandler)

	graphiqlHandler, err := handlers.NewGraphiQLHandler([]string{"graphiql"}, graphqlHandler.PathString())
	if err != nil {
		return nil, err
	}
	s.RegisterHandler(graphiqlHandler)

	return &s, nil
}
