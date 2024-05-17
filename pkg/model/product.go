package model

import (
	"Golang_Project/pkg/validator"
	"database/sql"
	"errors"
	"log"
	"sort"
	"time"
)

type Product struct {
	ID             string    `json:"id"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Price          float64   `json:"price"`
	AvailableStock uint      `json:"availableStock"`
	ShopID         int64     `json:"shopID"`
}

type ProductModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m *ProductModel) AddProduct(product Product) error {
	v := validator.New()

	v.Check(product.Title != "", "title", "Title is required")
	v.Check(product.Description != "", "description", "Description is required")
	v.Check(product.ShopID > 0, "shopID", "ShopID must be a positive integer")

	if !v.Valid() {
		return errors.New("invalid product data")
	}

	_, err := m.DB.Exec("INSERT INTO products (created_at, updated_at, title, description, price, shop_id) VALUES (NOW(), NOW(), $1, $2, $3, $4)",
		product.Title, product.Description, product.Price, product.ShopID)
	if err != nil {
		m.ErrorLog.Println("Error adding product:", err)
		return err
	}

	m.InfoLog.Println("Product added successfully")
	return nil
}

func (m *ProductModel) GetProductByID(id string) (*Product, error) {
	var product Product
	err := m.DB.QueryRow("SELECT id, created_at, updated_at, title, description, price, shop_id FROM products WHERE id = $1", id).
		Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Title, &product.Description, &product.Price, &product.ShopID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Product not found")
		}
		m.ErrorLog.Println("Error getting product:", err)
		return nil, err
	}
	return &product, nil
}

func (m *ProductModel) DeleteProductByID(id int) error {
	_, err := m.DB.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		m.ErrorLog.Println("Error deleting product:", err)
		return err
	}
	m.InfoLog.Println("Product deleted successfully")
	return nil
}

func (m *ProductModel) UpdateProductByID(id int, newData Product) error {
	v := validator.New()

	v.Check(newData.Title != "", "title", "Title is required")
	v.Check(newData.Description != "", "description", "Description is required")
	v.Check(newData.ShopID > 0, "shopID", "ShopID must be a positive integer")

	if !v.Valid() {
		return errors.New("invalid product data")
	}

	_, err := m.DB.Exec("UPDATE products SET title = $1, description = $2, price = $3, updated_at = $4 WHERE id = $5",
		newData.Title, newData.Description, newData.Price, time.Now(), id)
	if err != nil {
		m.ErrorLog.Println("Error updating product:", err)
		return err
	}
	m.InfoLog.Println("Product updated successfully")
	return nil
}

func (m *ProductModel) GetAllProduct(filters Filters) ([]Product, Metadata, error) {
	rows, err := m.DB.Query("SELECT id, created_at, updated_at, title, description, price, shop_id FROM products")
	if err != nil {
		m.ErrorLog.Println("Error getting products:", err)
		return nil, Metadata{}, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Title, &product.Description, &product.Price, &product.ShopID); err != nil {
			m.ErrorLog.Println("Error scanning product:", err)
			return nil, Metadata{}, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		m.ErrorLog.Println("Error iterating rows:", err)
		return nil, Metadata{}, err
	}

	if filters.Title != "" {
		products = FilterByTitle(products, filters.Title)
	}

	if filters.SortBy != "" {
		products = SortByPrice(products, filters.SortBy)
		if filters.SortBy == "price" && filters.SortOrder == "desc" {
			for i, j := 0, len(products)-1; i < j; i, j = i+1, j-1 {
				products[i], products[j] = products[j], products[i]
			}
		}
	}

	totalRecords := len(products)
	paginatedProducts := PaginateForProduct(products, filters.Page, filters.PageSize)
	metadata := CalculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return paginatedProducts, metadata, nil
}

func (m *ProductModel) SortProducts(products []Product, sortBy, sortOrder string) []Product {
	switch sortBy {
	case "price":
		sort.Slice(products, func(i, j int) bool {
			if sortOrder == "asc" {
				return products[i].Price < products[j].Price
			}
			return products[i].Price > products[j].Price
		})
	}

	return products
}
