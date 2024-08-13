package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/sosshik/users-service/internal/handlers"
	"github.com/sosshik/users-service/internal/repository"
	"github.com/sosshik/users-service/internal/service"
	"os"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	repos := repository.NewRepository()

	services := service.NewService(repos)

	handler := handlers.NewHandler(services)

	srv := handler.InitRoutes()
	log.Fatal(srv.Start(":8080"))
}
