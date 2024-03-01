package model

import (
	"database/sql"
	"log"
	"os"
)

type Models struct {
	Product ProductModel
	Shop    ShopModel
}

func NewModels(db *sql.DB) Models {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return Models{
		Product: ProductModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Shop: ShopModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		// Добавьте инициализацию других моделей, если необходимо
	}
}
