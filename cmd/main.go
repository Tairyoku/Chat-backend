package main

import (
	"cmd/pkg/handler"
	"cmd/pkg/repository"
	"cmd/pkg/service"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := repository.NewRepositoryDB(repository.Config{
		Username: os.Getenv("DBUsername"),
		Password: os.Getenv("DBPassword"),
		Host:     os.Getenv("DBHost"),
		Url:      os.Getenv("DBUrl"),
		DBName:   os.Getenv("DBName"),
	})
	if err != nil {
		log.Fatalf("error %s", err.Error())
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	server := new(service.Server)
	if err := server.Run(os.Getenv("PORT"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error %s", err.Error())
	}

}
