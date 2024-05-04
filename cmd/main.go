package main

import (
	"Golang_Project/api"
	"Golang_Project/pkg/model"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type config struct {
	port       int
	env        string
	fill       bool
	migrations string
	db         struct {
		dsn string
	}
}

func main() {
	fs := flag.NewFlagSet("demo-app", flag.ContinueOnError)

	var (
		cfg        config
		fill       = fs.Bool("fill", false, "Fill database with dummy data")
		migrations = fs.String("migrations", "", "Path to migration files folder. If not provided, migrations do not applied")
		port       = fs.Int("port", 8080, "API server port")
		env        = fs.String("env", "development", "Environment (development|staging|production)")
		dbDsn      = fs.String("dsn", "postgresql://postgres:Nurkhan05@db:5432/jana?sslmode=disable", "PostgreSQL DSN")
	)

	cfg.port = *port
	cfg.env = *env
	cfg.fill = *fill
	cfg.db.dsn = *dbDsn
	cfg.migrations = *migrations

	log.Println("starting application with configuration", map[string]string{
		"port":       fmt.Sprintf("%d", cfg.port),
		"fill":       fmt.Sprintf("%t", cfg.fill),
		"env":        cfg.env,
		"db":         cfg.db.dsn,
		"migrations": cfg.migrations,
	})

	// Установка соединения с базой данных
	db, err := openDB(cfg)
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

func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config // struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
