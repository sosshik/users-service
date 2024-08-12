package main

import (
	"github.com/sosshik/users-service/internal/handlers"
)

func main() {
	handler := handlers.NewHandler()
	srv := handler.InitRoutes()
	srv.Logger.Fatal(srv.Start(":8080"))
}
