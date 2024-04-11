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

//	func (m *ShopModel) GetAllShops(filters Filters) ([]Shop, Metadata, error) {
//		// Prepare SQL query
//		query := "SELECT id, created_at, updated_at, title, description, type FROM shop"
//
//		// Check if filtering criteria are provided
//		if filters.Type != "" {
//			query += " WHERE type = '" + filters.Type + "'"
//		}
//
//		// Execute the query against the database
//		rows, err := m.DB.Query(query)
//		if err != nil {
//			m.ErrorLog.Println("Error getting shops:", err)
//			return nil, Metadata{}, err
//		}
//		defer rows.Close()
//
//		// Scan the query results into Shop structure
//		var shops []Shop
//		for rows.Next() {
//			var shop Shop
//			if err := rows.Scan(&shop.Id, &shop.CreatedAt, &shop.UpdatedAt, &shop.Title, &shop.Description, &shop.Type); err != nil {
//				m.ErrorLog.Println("Error scanning shop:", err)
//				return nil, Metadata{}, err
//			}
//			shops = append(shops, shop)
//		}
//		if err := rows.Err(); err != nil {
//			m.ErrorLog.Println("Error iterating rows:", err)
//			return nil, Metadata{}, err
//		}
//
//		// Apply sorting
//		if filters.SortBy != "" {
//			switch filters.SortBy {
//			case "title":
//				sort.Slice(shops, func(i, j int) bool {
//					return shops[i].Title < shops[j].Title
//				})
//			// Add more cases for other sorting options if needed
//			default:
//				// Handle unknown sorting field
//				return nil, Metadata{}, errors.New("unknown sort field")
//			}
//		}
//
//		// Apply pagination
//		startIdx := (filters.Page - 1) * filters.PageSize
//		endIdx := startIdx + filters.PageSize
//		if endIdx > len(shops) {
//			endIdx = len(shops)
//		}
//		paginatedShops := shops[startIdx:endIdx]
//
//		// Calculate pagination metadata
//		totalRecords := len(shops)
//		metadata := CalculateMetadata(totalRecords, filters.Page, filters.PageSize)
//
//		return paginatedShops, metadata, nil
//	}
//
//	func (m *ShopModel) GetAllShops() ([]Shop, error) {
//		rows, err := m.DB.Query("SELECT id, created_at, updated_at, title, description,type FROM shop")
//		if err != nil {
//			m.ErrorLog.Println("Error getting shops:", err)
//			return nil, err
//		}
//		defer rows.Close()
//
//		var shops []Shop
//		for rows.Next() {
//			var shop Shop
//			if err := rows.Scan(&shop.Id, &shop.CreatedAt, &shop.UpdatedAt, &shop.Title, &shop.Description, &shop.Type); err != nil {
//				m.ErrorLog.Println("Error scanning shop:", err)
//				return nil, err
//			}
//			shops = append(shops, shop)
//		}
//		if err := rows.Err(); err != nil {
//			m.ErrorLog.Println("Error iterating rows:", err)
//			return nil, err
//		}
//
//		return shops, nil
//	}

func (m *ShopModel) GetAllShops(filters Filters) ([]Shop, Metadata, error) {
	// Fetch all shops from the database
	rows, err := m.DB.Query("SELECT id, created_at, updated_at, title, description, type FROM shop")
	if err != nil {
		m.ErrorLog.Println("Error getting shops:", err)
		return nil, Metadata{}, err
	}
	defer rows.Close()

	// Scan the query results into Shop structure
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

	// Apply filtering
	if filters.Type != "" {
		shops = FilterByType(shops, filters.Type)
	}

	// Apply sorting
	if filters.SortBy != "" {
		switch filters.SortBy {
		case "title":
			shops = SortByTitle(shops)
		// Add more cases for other sorting options if needed
		default:
			// Handle unknown sorting field
			return nil, Metadata{}, errors.New("unknown sort field")
		}
	}

	// Apply pagination
	paginatedShops := Paginate(shops, filters.Page, filters.PageSize)

	// Calculate pagination metadata
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
