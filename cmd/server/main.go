package main

import (
	"log"

	"github.com/Dimashey/gprc/internal/db"
	"github.com/Dimashey/gprc/internal/rocket"
	"github.com/Dimashey/gprc/internal/transport/grpc"
)

func Run() error {
	rocketStore, err := db.New()
	if err != nil {
		return err
	}
	err = rocketStore.Migrate()

	if err != nil {
		log.Println("Failed to run migrations")
		return err
	}

	rocketService := rocket.New(rocketStore)

	rktHandler := grpc.New(rocketService)

	if err := rktHandler.Serve(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}
