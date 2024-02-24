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

var shops = []Shop{
	{
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Location    string    `json:"location"`
		Owner       *User     `json:"owner"`

	},
	{
		ID          string  `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Quantity    int     `json:"quantity"`
		Category    string  `json:"category"`
		Remainder   int      `json:"remainder"`
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
