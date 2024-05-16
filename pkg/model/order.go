package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"
)

type Order struct {
	ID           string `json:"id"`
	UserID       int
	Products     []*Product
	TotalAmount  float64
	DeliveryAddr string
	Status       string
	CreatedAt    time.Time
}

type OrderModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (o *Order) CalculateTotal() float64 {
	total := 0.0
	for _, product := range o.Products {
		total += product.Price
	}
	o.TotalAmount = total
	return total
}

func (m *OrderModel) CreateOrder(order *Order) error {
	order.CalculateTotal()

	productsJSON, err := json.Marshal(order.Products)
	if err != nil {
		m.ErrorLog.Println("Error marshalling order products:", err)
		return err
	}

	_, err = m.DB.Exec("INSERT INTO orders (user_id, products, total_amount, delivery_addr, status, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		order.UserID, productsJSON, order.TotalAmount, order.DeliveryAddr, order.Status, order.CreatedAt)
	if err != nil {
		m.ErrorLog.Println("Error creating order:", err)
		return err
	}

	m.InfoLog.Println("Order created successfully")
	return nil
}

func (m *OrderModel) UpdateOrderStatus(orderID string, status string) error {
	_, err := m.DB.Exec("UPDATE orders SET status = $1 WHERE id = $2", status, orderID)
	if err != nil {
		m.ErrorLog.Println("Error updating order status:", err)
		return err
	}

	m.InfoLog.Println("Order status updated successfully")
	return nil
}

//func (m *OrderModel) GetOrder(orderID string) (*Order, error) {
//	row := m.DB.QueryRow("SELECT id, user_id, products, total_amount, delivery_addr, status, created_at FROM orders WHERE id = $1", orderID)
//
//	order := &Order{}
//	var productsJSON []byte
//
//	err := row.Scan(&order.ID, &order.UserID, &productsJSON, &order.TotalAmount, &order.DeliveryAddr, &order.Status, &order.CreatedAt)
//	if err != nil {
//		m.ErrorLog.Println("Error getting order:", err)
//		return nil, err
//	}
//
//	err = json.Unmarshal(productsJSON, &order.Products)
//	if err != nil {
//		m.ErrorLog.Println("Error unmarshalling order products:", err)
//		return nil, err
//	}
//
//	return order, nil
//}

func (m *OrderModel) GetOrder(orderID string) (*Order, error) {
	row := m.DB.QueryRow("SELECT id, user_id, products, total_amount, delivery_addr, status, created_at FROM orders WHERE id = $1", orderID)

	order := &Order{}
	var productsJSON []byte

	err := row.Scan(&order.ID, &order.UserID, &productsJSON, &order.TotalAmount, &order.DeliveryAddr, &order.Status, &order.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Order not found")
		}
		m.ErrorLog.Println("Error getting order:", err)
		return nil, err
	}

	err = json.Unmarshal(productsJSON, &order.Products)
	if err != nil {
		m.ErrorLog.Println("Error unmarshalling order products:", err)
		return nil, err
	}

	return order, nil
}

func (m *OrderModel) GetAllOrders(userID int) ([]*Order, error) {
	rows, err := m.DB.Query("SELECT id, user_id, products, total_amount, delivery_addr, status, created_at FROM orders WHERE user_id = $1", userID)
	if err != nil {
		m.ErrorLog.Println("Error getting orders:", err)
		return nil, err
	}
	defer rows.Close()

	orders := []*Order{}
	for rows.Next() {
		order := &Order{}
		var productsJSON []byte

		err := rows.Scan(&order.ID, &order.UserID, &productsJSON, &order.TotalAmount, &order.DeliveryAddr, &order.Status, &order.CreatedAt)
		if err != nil {
			m.ErrorLog.Println("Error scanning order:", err)
			return nil, err
		}

		err = json.Unmarshal(productsJSON, &order.Products)
		if err != nil {
			m.ErrorLog.Println("Error unmarshalling order products:", err)
			return nil, err
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		m.ErrorLog.Println("Error iterating rows:", err)
		return nil, err
	}

	return orders, nil
}
