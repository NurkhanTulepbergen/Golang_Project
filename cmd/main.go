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

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Bayernmunichtm25"
	dbname   = "jana"
)

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
	userModel := &model.UserModel{
		DB:       db,
		InfoLog:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		ErrorLog: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
	tokenModel := &model.TokenModel{
		DB:       db,
		InfoLog:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		ErrorLog: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
	permissionModel := &model.PermissionModel{
		DB:       db,
		InfoLog:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		ErrorLog: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}

	// Создание объекта Filters с необходимыми значениями
	filters := model.Filters{
		Page:     1,
		PageSize: 10,
		// Можете также установить другие значения фильтрации здесь, если необходимо
	}

	// Получение списка магазинов с помощью метода GetShops
	shops, metadata, err := shopModel.GetAllShops(filters)
	if err != nil {
		log.Println("Error getting shops:", err)
		return
	}

	// Вывод информации о магазинах и метаданных пагинации
	log.Println("Shops:", shops)
	log.Println("Metadata:", metadata)
	//log.Println("Users:", users)

	api := api.NewAPI(shopModel, productModel, userModel, tokenModel, permissionModel)
	api.StartServer()
}
