package model

import (
	"database/sql"
	"errors"
	"log"
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
}

type ProductModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

//	func (m *ProductModel) Insert(product *Product) error {
//		query := `
//	       INSERT INTO products (title, description, price, available_stock, created_at, updated_at)
//	       VALUES ($1, $2, $3, $4, $5, $6)
//	       RETURNING id, created_at, updated_at
//	       `
//		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//		defer cancel()
//
//		err := m.DB.QueryRowContext(ctx, query, product.Title, product.Description, product.Price, product.AvailableStock, time.Now(), time.Now()).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
//		if err != nil {
//			m.ErrorLog.Println("Error inserting product:", err)
//			return err
//		}
//
//		return nil
//	}
func (m *ProductModel) AddProduct(product Product) error {
	// Check if the shop data is valid
	if product.Title == "" || product.Description == "" {
		return errors.New("title, description, price are required fields")
	}

	// Perform the database insertion
	_, err := m.DB.Exec("INSERT INTO products (created_at, updated_at, title, description,price) VALUES (NOW(), NOW(), $1, $2, $3)",
		product.Title, product.Description, product.Price)
	if err != nil {
		m.ErrorLog.Println("Error adding product:", err)
		return err
	}

	m.InfoLog.Println("Product added successfully")
	return nil
}

//func (m *ProductModel) Get(id string) (*Product, error) {
//	query := `
//        SELECT id, created_at, updated_at, title, description, price, available_stock
//        FROM products
//        WHERE id = $1
//        `
//	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//	defer cancel()
//
//	var product Product
//	err := m.DB.QueryRowContext(ctx, query, id).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Title, &product.Description, &product.Price, &product.AvailableStock)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			return nil, nil
//		}
//		m.ErrorLog.Println("Error getting product:", err)
//		return nil, err
//	}
//
//	return &product, nil
//}

func (m *ProductModel) GetProductByID(id string) (*Product, error) {
	var product Product
	err := m.DB.QueryRow("SELECT id, created_at, updated_at, title, description,price FROM products WHERE id = $1", id).
		Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Title, &product.Description, &product.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Product not found")
		}
		m.ErrorLog.Println("Error getting product:", err)
		return nil, err
	}
	return &product, nil
}

//func (m *ProductModel) Update(product *Product) error {
//	query := `
//        UPDATE products
//        SET title = $1, description = $2, price = $3, available_stock = $4, updated_at = $5
//        WHERE id = $6
//        RETURNING updated_at
//        `
//	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//	defer cancel()
//
//	err := m.DB.QueryRowContext(ctx, query, product.Title, product.Description, product.Price, product.AvailableStock, time.Now(), product.ID).Scan(&product.UpdatedAt)
//	if err != nil {
//		m.ErrorLog.Println("Error updating product:", err)
//		return err
//	}
//
//	return nil
//}

//func (m *ProductModel) Delete(id string) error {
//	query := `
//        DELETE FROM products
//        WHERE id = $1
//        `
//	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//	defer cancel()
//
//	_, err := m.DB.ExecContext(ctx, query, id)
//	if err != nil {
//		m.ErrorLog.Println("Error deleting product:", err)
//		return err
//	}
//
//	return nil
//}

func (m *ProductModel) DeleteProductByID(id int) error {
	// Выполняем SQL-запрос для удаления магазина по его ID
	_, err := m.DB.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		m.ErrorLog.Println("Error deleting product:", err)
		return err
	}
	m.InfoLog.Println("Product deleted successfully")
	return nil
}

func (m *ProductModel) UpdateProductByID(id int, newData Product) error {
	_, err := m.DB.Exec("UPDATE products SET title = $1, description = $2, price = $3, updated_at = $4 WHERE id = $5",
		newData.Title, newData.Description, newData.Price, time.Now(), id)
	if err != nil {
		m.ErrorLog.Println("Error updating product:", err)
		return err
	}
	m.InfoLog.Println("Product updated successfully")
	return nil
}

func (m *ProductModel) GetAllProduct() ([]Product, error) {
	rows, err := m.DB.Query("SELECT id, created_at, updated_at, title, description, price FROM products")
	if err != nil {
		m.ErrorLog.Println("Error getting product:", err)
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt, &product.Title, &product.Description, &product.Price); err != nil {
			m.ErrorLog.Println("Error scanning product:", err)
			return nil, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		m.ErrorLog.Println("Error iterating rows:", err)
		return nil, err
	}

	return products, nil
}
