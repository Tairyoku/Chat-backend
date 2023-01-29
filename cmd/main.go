package main

import (
	"cmd/pkg/handler"
	"cmd/pkg/repository"
	"cmd/pkg/service"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	//"gorm.io/gorm"
	"log"
	"os"
)

func GetConnectionString() string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "3307"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "root"
	}

	password := os.Getenv("DB_PASS")
	if password == "" {
		password = "@root"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "chatDB"
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)
}
func main() {

	go handler.Hub.Run()
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Error loading .env file")
	}
	var db, err = gorm.Open("mysql", GetConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	//db, err := repository.NewRepositoryDB(repository.Config{
	//	Username: os.Getenv("DBUsername"),
	//	Password: os.Getenv("DBPassword"),
	//	Host:     os.Getenv("DBHost"),
	//	Url:      os.Getenv("DBUrl"),
	//	DBName:   os.Getenv("DBName"),
	//})
	//if err != nil {
	//	log.Fatalf("error %s", err.Error())
	//}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	server := new(service.Server)

	if err := server.Run(os.Getenv("PORT"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error %s", err.Error())
	}

}
