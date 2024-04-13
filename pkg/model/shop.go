package model

import (
	"database/sql"
	"errors"
	"log"
	"time"
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
	// Check if the shop data is valid
	if shop.Title == "" || shop.Description == "" {
		return errors.New("title and description are required fields")
	}

	// Perform the database insertion
	_, err := m.DB.Exec("INSERT INTO shop (created_at, updated_at, title, description, type) VALUES (NOW(), NOW(), $1, $2, $3)",
		shop.Title, shop.Description, shop.Type)
	if err != nil {
		m.ErrorLog.Println("Error adding shop:", err)
		return err
	}

	m.InfoLog.Println("Shop added successfully")
	return nil
}

func (m *ShopModel) UpdateShopByID(id int, newData Shop) error {
	_, err := m.DB.Exec("UPDATE shop SET title = $1, description = $2, type = $3, updated_at = $4 WHERE id = $5",
		newData.Title, newData.Description, newData.Type, time.Now(), id)
	if err != nil {
		m.ErrorLog.Println("Error updating shop:", err)
		return err
	}
	m.InfoLog.Println("Shop updated successfully")
	return nil
}

func (m *ShopModel) DeleteShopByID(id int) error {
	// Execute SQL query to delete a shop by its ID
	_, err := m.DB.Exec("DELETE FROM shop WHERE id = $1", id)
	if err != nil {
		m.ErrorLog.Println("Error deleting shop:", err)
		return err
	}
	m.InfoLog.Println("Shop deleted successfully")
	return nil
}

//func (m *ShopModel) GetAllShopsWithFilters(filters map[string]string, sortBy string, sortOrder string, page int, limit int) ([]Shop, error) {
//	// Формируем SQL-запрос с учетом фильтров, сортировки и пагинации
//	query := "SELECT id, created_at, updated_at, title, description, type FROM shop"
//
//	// Формируем условие для фильтрации
//	var args []interface{}
//	if len(filters) > 0 {
//		query += " WHERE "
//		for key, value := range filters {
//			query += key + " = ? AND "
//			args = append(args, value)
//		}
//		// Убираем последний "AND"
//		query = query[:len(query)-5]
//	}
//
//	// Добавляем сортировку
//	if sortBy != "" {
//		query += " ORDER BY " + sortBy
//		if sortOrder != "" {
//			query += " " + sortOrder
//		}
//	}
//
//	// Добавляем пагинацию
//	if limit > 0 {
//		query += " LIMIT ? OFFSET ?"
//		args = append(args, limit, (page-1)*limit)
//	}
//
//	// Выполняем запрос
//	rows, err := m.DB.Query(query, args...)
//	if err != nil {
//		m.ErrorLog.Println("Error getting shops:", err)
//		return nil, err
//	}
//	defer rows.Close()
//
//	// Считываем результаты запроса
//	var shops []Shop
//	for rows.Next() {
//		var shop Shop
//		if err := rows.Scan(&shop.Id, &shop.CreatedAt, &shop.UpdatedAt, &shop.Title, &shop.Description, &shop.Type); err != nil {
//			m.ErrorLog.Println("Error scanning shop:", err)
//			return nil, err
//		}
//		shops = append(shops, shop)
//	}
//	if err := rows.Err(); err != nil {
//		m.ErrorLog.Println("Error iterating rows:", err)
//		return nil, err
//	}
//
//	return shops, nil
//}

//func (m *ShopModel) GetAllShops() ([]Shop, error) {
//	rows, err := m.DB.Query("SELECT id, created_at, updated_at, title, description,type FROM shop")
//	if err != nil {
//		m.ErrorLog.Println("Error getting shops:", err)
//		return nil, err
//	}
//	defer rows.Close()
//
//	var shops []Shop
//	for rows.Next() {
//		var shop Shop
//		if err := rows.Scan(&shop.Id, &shop.CreatedAt, &shop.UpdatedAt, &shop.Title, &shop.Description, &shop.Type); err != nil {
//			m.ErrorLog.Println("Error scanning shop:", err)
//			return nil, err
//		}
//		shops = append(shops, shop)
//	}
//	if err := rows.Err(); err != nil {
//		m.ErrorLog.Println("Error iterating rows:", err)
//		return nil, err
//	}
//
//	return shops, nil
//}

func (m *ShopModel) GetAllShops(filters Filters) ([]Shop, Metadata, error) {
	// Fetch all shops from the database
	rows, err := m.DB.Query("SELECT id, created_at, updated_at, title, description, type FROM shop")
	if err != nil {
		m.ErrorLog.Println("Error getting shops:", err)
		return nil, Metadata{}, err
	}
	defer rows.Close()

	var shops []Shop
	for rows.Next() {
		var shop Shop
		if err := rows.Scan(&shop.Id, &shop.CreatedAt, &shop.UpdatedAt, &shop.Title, &shop.Description, &shop.Type); err != nil {
			m.ErrorLog.Println("Error scanning shop:", err)
			return nil, Metadata{}, err
		}
		shops = append(shops, shop)
	}
	if err := rows.Err(); err != nil {
		m.ErrorLog.Println("Error iterating rows:", err)
		return nil, Metadata{}, err
	}

	if filters.Type != "" {
		shops = FilterByType(shops, filters.Type)
	}

	if filters.SortBy != "" {
		switch filters.SortBy {
		case "title":
			shops = SortByTitle(shops)

		default:

			return nil, Metadata{}, errors.New("unknown sort field")
		}
	}

	paginatedShops := Paginate(shops, filters.Page, filters.PageSize)

	totalRecords := len(shops)
	metadata := CalculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return paginatedShops, metadata, nil
}

func (m *ShopModel) GetShopByID(id string) (*Shop, error) {
	var shop Shop
	err := m.DB.QueryRow("SELECT id, created_at, updated_at, title, description, type FROM shop WHERE id = $1", id).
		Scan(&shop.Id, &shop.CreatedAt, &shop.UpdatedAt, &shop.Title, &shop.Description, &shop.Type)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Shop not found")
		}
		m.ErrorLog.Println("Error getting shop:", err)
		return nil, err
	}
	return &shop, nil
}
