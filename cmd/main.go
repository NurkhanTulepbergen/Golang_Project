package main

import (
	"Golang_Project/api"
	"Golang_Project/pkg/model"
	"database/sql"
	"flag"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
)

// Define the configuration struct
type config struct {
	port       int
	env        string
	migrations string
	db         struct {
		dsn string
	}
}

func main() {
	fs := flag.NewFlagSet("demo-app", flag.ContinueOnError)

	var (
		cfg        config
		migrations = fs.String("migrations", "", "Path to migration files folder. If not provided, migrations do not applied")
		port       = fs.Int("port", 8080, "API server port")
		env        = fs.String("env", "development", "Environment (development|staging|production)")
		dbDsn      = fs.String("dsn", "postgres://postgres:Bayernmunichtm25@localhost:5432/jana?sslmode=disable", "PostgreSQL DSN")
	)

	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatalf("error parsing flags: %v", err)
	}

	cfg.port = *port
	cfg.env = *env
	cfg.db.dsn = *dbDsn
	cfg.migrations = *migrations

	// Initialize logger
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Connect to database
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		logger.Fatalf("error connecting to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if cfg.migrations != "" {
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			logger.Fatalf("error creating migration driver: %v", err)
		}
		m, err := migrate.NewWithDatabaseInstance(
			cfg.migrations,
			"postgres", driver)
		if err != nil {
			logger.Fatalf("error creating migration instance: %v", err)
		}
		if err := m.Up(); err != nil {
			logger.Fatalf("error applying migrations: %v", err)
		}
	}

	// Create models
	shopModel := &model.ShopModel{
		DB:       db,
		InfoLog:  logger,
		ErrorLog: logger,
	}
	productModel := &model.ProductModel{
		DB:       db,
		InfoLog:  logger,
		ErrorLog: logger,
	}
	userModel := &model.UserModel{
		DB:       db,
		InfoLog:  logger,
		ErrorLog: logger,
	}
	tokenModel := &model.TokenModel{
		DB:       db,
		InfoLog:  logger,
		ErrorLog: logger,
	}
	permissionModel := &model.PermissionModel{
		DB:       db,
		InfoLog:  logger,
		ErrorLog: logger,
	}
	cartModel := &model.CartModel{
		DB:       db,
		InfoLog:  logger,
		ErrorLog: logger,
	}
	orderModel := &model.OrderModel{
		DB:       db,
		InfoLog:  logger,
		ErrorLog: logger,
	}
	historyModel := &model.HistoryModel{
		DB:       db,
		InfoLog:  logger,
		ErrorLog: logger,
	}
	followModel := &model.FollowModel{
		DB:       db,
		InfoLog:  logger,
		ErrorLog: logger,
	}

	// Start server
	api := api.NewAPI(shopModel, productModel, userModel, tokenModel, permissionModel, cartModel, orderModel, historyModel, followModel)
	api.StartServer(cfg.port)
}
