package model

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

type OrderProduct struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type Order struct {
	ID           int            `json:"id"`
	UserID       int            `json:"user_id"`
	Products     []OrderProduct `json:"products"`
	TotalAmount  float64        `json:"total_amount"`
	DeliveryAddr string         `json:"delivery_addr"`
	Status       string         `json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
}

type OrderModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

// Calculate the total amount for an order
func (o *Order) CalculateTotal(db *sql.DB) (float64, error) {
	total := 0.0
	for _, op := range o.Products {
		// Fetch product price from the database using ProductID
		var price float64
		err := db.QueryRow("SELECT price FROM products WHERE id = $1", op.ProductID).Scan(&price)
		if err != nil {
			return 0, err
		}
		total += price * float64(op.Quantity)
	}
	o.TotalAmount = total
	log.Println("Calculated total amount:", total) // Add this line to log the total amount
	return total, nil
}

func (m *OrderModel) CreateOrder(order *Order) error {
	// Calculate the total amount for the order
	total, err := order.CalculateTotal(m.DB)
	if err != nil {
		m.ErrorLog.Println("Error calculating total:", err)
		return err
	}

	// Start a transaction
	tx, err := m.DB.Begin()
	if err != nil {
		m.ErrorLog.Println("Error starting transaction:", err)
		return err
	}
	defer tx.Rollback()

	order.CreatedAt = time.Now()

	// Insert the order
	err = tx.QueryRow(
		"INSERT INTO orders (user_id, delivery_addr, status, total_amount, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		order.UserID, order.DeliveryAddr, order.Status, total, order.CreatedAt,
	).Scan(&order.ID)

	if err != nil {
		m.ErrorLog.Println("Error inserting order:", err)
		return err
	}

	// Insert the order products
	for _, op := range order.Products {
		_, err := tx.Exec(
			"INSERT INTO order_products (order_id, product_id, quantity) VALUES ($1, $2, $3)",
			order.ID, op.ProductID, op.Quantity,
		)
		if err != nil {
			m.ErrorLog.Println("Error inserting order product:", err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		m.ErrorLog.Println("Error committing transaction:", err)
		return err
	}

	m.InfoLog.Println("Order created successfully:", order.ID)
	return nil
}

// Fetch all orders for a specific user

func (m *OrderModel) GetAllOrders(userID int) ([]*Order, error) {
	rows, err := m.DB.Query("SELECT id, user_id, total_amount, delivery_addr, status, created_at FROM orders WHERE user_id = $1", userID)
	if err != nil {
		m.ErrorLog.Println("Error getting orders:", err)
		return nil, err
	}
	defer rows.Close()

	var orders []*Order
	for rows.Next() {
		order := &Order{}
		err := rows.Scan(&order.ID, &order.UserID, &order.TotalAmount, &order.DeliveryAddr, &order.Status, &order.CreatedAt)
		if err != nil {
			m.ErrorLog.Println("Error scanning order:", err)
			return nil, err
		}

		// Fetch products for this order
		prodRows, err := m.DB.Query("SELECT product_id, quantity FROM order_products WHERE order_id = $1", order.ID)
		if err != nil {
			m.ErrorLog.Println("Error getting order products:", err)
			return nil, err
		}
		defer prodRows.Close()

		for prodRows.Next() {
			op := OrderProduct{}
			err := prodRows.Scan(&op.ProductID, &op.Quantity)
			if err != nil {
				m.ErrorLog.Println("Error scanning order product:", err)
				return nil, err
			}
			order.Products = append(order.Products, op)
		}

		if err := prodRows.Err(); err != nil {
			m.ErrorLog.Println("Error iterating over product rows:", err)
			return nil, err
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		m.ErrorLog.Println("Error iterating over order rows:", err)
		return nil, err
	}

	return orders, nil
}

func (m *OrderModel) GetOrder(orderID int) (*Order, error) {
	order := &Order{}
	err := m.DB.QueryRow("SELECT id, user_id, total_amount, delivery_addr, status, created_at FROM orders WHERE id = $1",
		orderID).Scan(&order.ID, &order.UserID, &order.TotalAmount, &order.DeliveryAddr, &order.Status, &order.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order not found")
		}
		m.ErrorLog.Println("Error getting order:", err)
		return nil, err
	}

	// Fetch products for this order
	prodRows, err := m.DB.Query("SELECT product_id, quantity FROM order_products WHERE order_id = $1", order.ID)
	if err != nil {
		m.ErrorLog.Println("Error getting order products:", err)
		return nil, err
	}
	defer prodRows.Close()

	for prodRows.Next() {
		op := OrderProduct{}
		err := prodRows.Scan(&op.ProductID, &op.Quantity)
		if err != nil {
			m.ErrorLog.Println("Error scanning order product:", err)
			return nil, err
		}
		order.Products = append(order.Products, op)
	}

	if err := prodRows.Err(); err != nil {
		m.ErrorLog.Println("Error iterating over product rows:", err)
		return nil, err
	}

	return order, nil
}
func (m *OrderModel) FilterOrders(userID int, filters Filters) ([]*Order, Metadata, error) {
	// Fetch all orders for a specific user
	orders, err := m.GetAllOrders(userID)
	if err != nil {
		return nil, Metadata{}, err
	}

	// Apply filters
	if filters.Title != "" {
		orders = FilterByOrder(orders, filters.Title)
	}

	// Apply sorting
	if filters.SortBy != "" {
		switch filters.SortBy {
		case "total_amount":
			orders = SortByTotalAmount(orders, filters.SortOrder)
		case "created_at":
			orders = SortByCreatedAt(orders, filters.SortOrder)
			// Add more cases for additional fields if needed
		}
	}

	// Calculate metadata
	totalRecords := len(orders)
	metadata := CalculateMetadata(totalRecords, filters.Page, filters.PageSize)

	// Apply pagination
	orders = PaginateOrders(orders, filters.Page, filters.PageSize)

	return orders, metadata, nil
}

func (m *OrderModel) DeleteOrder(orderID int) error {
	// Start a transaction
	tx, err := m.DB.Begin()
	if err != nil {
		m.ErrorLog.Println("Error starting transaction:", err)
		return err
	}
	defer tx.Rollback()

	// Delete order products associated with the order
	_, err = tx.Exec("DELETE FROM order_products WHERE order_id = $1", orderID)
	if err != nil {
		m.ErrorLog.Println("Error deleting order products:", err)
		return err
	}

	// Delete the order itself
	_, err = tx.Exec("DELETE FROM orders WHERE id = $1", orderID)
	if err != nil {
		m.ErrorLog.Println("Error deleting order:", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		m.ErrorLog.Println("Error committing transaction:", err)
		return err
	}

	m.InfoLog.Println("Order deleted successfully:", orderID)
	return nil
}
