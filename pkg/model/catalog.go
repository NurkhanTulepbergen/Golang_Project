package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Product struct {
	ID          int     `json:"id"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

type ProductModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m ProductModel) Insert(product *Product) error {
	// Insert a new product into the database.
	query := `
		INSERT INTO products (name, description, price, quantity) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at, updated_at
		`
	args := []interface{}{product.Name, product.Description, product.Price, product.Quantity}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
}

func (m ProductModel) Get(id int) (*Product, error) {
	// Retrieve a specific product based on its ID.
	query := `
		SELECT id, created_at, updated_at, name, description, price, quantity
		FROM products
		WHERE id = $1
		`
	var product Product
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Name, &product.Description, &product.Price, &product.Quantity)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (m ProductModel) Update(product *Product) error {
	// Update a specific product in the database.
	query := `
		UPDATE products
		SET name = $1, description = $2, price = $3, quantity = $4
		WHERE id = $5
		RETURNING updated_at
		`
	args := []interface{}{product.Name, product.Description, product.Price, product.Quantity, product.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&product.UpdatedAt)
}

func (m ProductModel) Delete(id int) error {
	// Delete a specific product from the database.
	query := `
		DELETE FROM products
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
