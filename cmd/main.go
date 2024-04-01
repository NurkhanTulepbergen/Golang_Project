package main

import (
	"Golang_Project/api"
	"Golang_Project/pkg/model"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

//
//const (
//	host     = "localhost"
//	port     = 5432
//	user     = "postgres"
//	password = "adminkbtu"
//	dbname   = "jana"
//)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Установка соединения с базой данных
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	shopModel := &model.ShopModel{
		DB:       db,
		InfoLog:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		ErrorLog: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
	productModel := &model.ProductModel{
		DB:       db,
		InfoLog:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		ErrorLog: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}

	api := api.NewAPI(shopModel, productModel)
	api.StartServer()
}
