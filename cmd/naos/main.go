package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/OpenPeeDeeP/xdg"
	"gitlab.com/Dophin2009/nao/pkg/data"
)

func main() {
	// Exit with status code 0 at the end
	defer os.Exit(0)

	println("-------------------: NAO SERVER :-------------------")

	// Read configuration files
	etcDir := "/etc/nao/"
	userDir := xdg.ConfigHome() + "/nao/"
	confFileDirs := []string{etcDir, userDir}
	conf, err := ReadConfig(confFileDirs)
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
	serverAddress := fmt.Sprintf("%s:%s", conf.Hostname, conf.Port)
	controller := ControllerNew(db)
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
