package model

import (
	"database/sql"
	"errors"
	"log"
	"os"
)

var (

	// ErrEditConflict is returned when a there is a data race, and we have an edit conflict.
	ErrEditConflict   = errors.New("edit conflict")
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Product ProductModel
	Shop    ShopModel
	//Cart       CartModel
	User       UserModel
	Token      TokenModel
	Permission PermissionModel
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
		User: UserModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Token: TokenModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Permission: PermissionModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		// Добавьте инициализацию других моделей, если необходимо
	}
}
