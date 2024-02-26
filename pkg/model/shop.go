package model

import (
	"database/sql"
	"errors"
	"log"
)

type Shop struct {
	Id          string `json:"id"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type ShopModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

type Tmp struct {
	Shop
	ShopModel
}

var shops = []Shop{
	{
		
	},
	{
		
	},
	// Add other shops as needed
}

func GetShops() []Shop {
	return shops
}

func GetShop(id string) (*Shop, error) {
	for _, s := range shops {
		if s.Id == id {
			return &s, nil
		}
	}
	return nil, errors.New("Shop not found")
}
