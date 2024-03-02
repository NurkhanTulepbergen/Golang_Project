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

func (m *ShopModel) AddShop(shop Shop) error {
	_, err := m.DB.Exec("INSERT INTO shop (created_at, updated_at, title, description) VALUES (NOW(), NOW(), $1, $2)",
		shop.Title, shop.Description)
	if err != nil {
		m.ErrorLog.Println("Error adding shop:", err)
		return err
	}
	m.InfoLog.Println("Shop added successfully")
	return nil
}
func (m *ShopModel) GetShopByID(id string) (*Shop, error) {
	var shop Shop
	err := m.DB.QueryRow("SELECT id, created_at, updated_at, title, description FROM shop WHERE id = $1", id).
		Scan(&shop.Id, &shop.CreatedAt, &shop.UpdatedAt, &shop.Title, &shop.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Shop not found")
		}
		m.ErrorLog.Println("Error getting shop:", err)
		return nil, err
	}
	return &shop, nil
}
func (m *ShopModel) DeleteShopByID(id int) error {
	// Выполняем SQL-запрос для удаления магазина по его ID
	_, err := m.DB.Exec("DELETE FROM shop WHERE id = $1", id)
	if err != nil {
		m.ErrorLog.Println("Error deleting shop:", err)
		return err
	}
	m.InfoLog.Println("Shop deleted successfully")
	return nil
}
func (m *ShopModel) GetAllShops() ([]Shop, error) {
	rows, err := m.DB.Query("SELECT id, created_at, updated_at, title, description FROM shop")
	if err != nil {
		m.ErrorLog.Println("Error getting shops:", err)
		return nil, err
	}
	defer rows.Close()

	var shops []Shop
	for rows.Next() {
		var shop Shop
		if err := rows.Scan(&shop.Id, &shop.CreatedAt, &shop.UpdatedAt, &shop.Title, &shop.Description); err != nil {
			m.ErrorLog.Println("Error scanning shop:", err)
			return nil, err
		}
		shops = append(shops, shop)
	}
	if err := rows.Err(); err != nil {
		m.ErrorLog.Println("Error iterating rows:", err)
		return nil, err
	}

	return shops, nil
}
