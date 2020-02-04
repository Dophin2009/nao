package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/Dophin2009/nao/internal/naos"
)

// TODO: Parse command line flags

func main() {
	// Exit with status code 0 at the end
	defer os.Exit(0)

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	// Read configuration files
	conf, err := naos.ReadConfigs()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
		return
	}

	s, err := naos.NewApplication(conf)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
		return
	}
	defer s.DataLayer.Database.Close()

	// Launch server in goroutine
	shttp := s.HTTPServer()
	go func() {
		log.WithFields(log.Fields{
			"address": shttp.Addr,
		}).Info("Launching server")
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
