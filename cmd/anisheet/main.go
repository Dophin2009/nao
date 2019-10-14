package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"gitlab.com/Dophin2009/anisheet/cmd/anisheet/controller"
	"gitlab.com/Dophin2009/anisheet/pkg/data"
	bolt "go.etcd.io/bbolt"
)

var db *bolt.DB

func main() {
	// Exit with status code 0 at the end
	defer os.Exit(0)

	// Open database connection
	db, err := data.ConnectDatabase("/tmp/anisheet.db", true)
	if err != nil {
		panic("error connecting to database ")
	}
	// Clear database and close connection at the end
	defer db.Close()
	defer data.ClearDatabase(db)

	// Create the API controller and HTTP server
	controller := controller.NewController(db)
	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: controller.Router,
	}

	// Allow user to set graceful timeout duration
	var wait time.Duration
	flag.DurationVar(&wait, "gt", time.Second*15, "graceful timeout: the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// Launch server in goroutine
	go func() {
		// err := server.ListenAndServeTLS("cert.pem", "key.pem")
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for SIGINTERRUPT signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Wait for processes to end, then shutdown
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	server.Shutdown(ctx)

	println()
	log.Println("Exiting...")
}
