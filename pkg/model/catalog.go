package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Product struct {
	ID             string  `json:"id"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      string  `json:"updatedAt"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	Price          float64 `json:"price"`
	AvailableStock uint    `json:"availableStock"`
}

type ProductModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m ProductModel) Insert(product *Product) error {
	query := `
		INSERT INTO products (title, description, price, available_stock, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id, created_at, updated_at
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query and scan the result into the product struct
	err := m.DB.QueryRowContext(ctx, query, product.Title, product.Description, product.Price, product.AvailableStock, time.Now(), time.Now()).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		m.ErrorLog.Println("Error inserting product:", err)
		return err
	}

	return nil
}

func (m ProductModel) Get(id string) (*Product, error) {
	query := `
		SELECT id, created_at, updated_at, title, description, price, available_stock
		FROM products
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var product Product
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Title, &product.Description, &product.Price, &product.AvailableStock)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil if the product is not found
		}
		m.ErrorLog.Println("Error getting product:", err)
		return nil, err
	}

	return &product, nil
}

func (m ProductModel) Update(product *Product) error {
	query := `
		UPDATE products
		SET title = $1, description = $2, price = $3, available_stock = $4, updated_at = $5
		WHERE id = $6
		RETURNING updated_at
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query and scan the result into the product struct
	err := m.DB.QueryRowContext(ctx, query, product.Title, product.Description, product.Price, product.AvailableStock, time.Now(), product.ID).Scan(&product.UpdatedAt)
	if err != nil {
		m.ErrorLog.Println("Error updating product:", err)
		return err
	}

	return nil
}

func (m ProductModel) Delete(id string) error {
	query := `
		DELETE FROM products
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query to delete the product
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		m.ErrorLog.Println("Error deleting product:", err)
		return err
	}

	return nil
}
