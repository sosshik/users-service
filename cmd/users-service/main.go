package main

import (
	log "github.com/sirupsen/logrus"
	_ "github.com/sosshik/users-service/docs"
	"github.com/sosshik/users-service/internal/handlers"
	"github.com/sosshik/users-service/internal/repository"
	"github.com/sosshik/users-service/internal/service"
	"os"
)

// @title Users Service API
// @version 1.0
// @description This is a sample service for managing users.

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8090
// @BasePath /

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	repos := repository.NewRepository()

	services := service.NewService(repos)

	handler := handlers.NewHandler(services)

	srv := handler.InitRoutes()
	log.Fatal(srv.Start(":8090"))
}
