package main

import (
	"context"
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
	println()
	log.Println("Establishing database connection")
	db, err := data.ConnectDatabase("/tmp/anisheet.db", true)
	if err != nil {
		log.Fatal("Error connecting to database ")
		return
	}
	// Clear database and close connection at the end
	defer db.Close()
	defer data.ClearDatabase(db)

	// Create the API controller and HTTP server
	const serverAddress = "0.0.0.0:8080"
	controller := controller.New(db)
	server := &http.Server{
		Addr:    serverAddress,
		Handler: controller.Router,
	}

	// Launch server in goroutine
	go func() {
		log.Println("Launching server on", serverAddress)
		// err := server.ListenAndServeTLS("cert.pem", "key.pem")
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for SIGINTERRUPT signal
	wait := time.Second * 15
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
