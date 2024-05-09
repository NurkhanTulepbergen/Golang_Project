package model

//
//import (
//	"database/sql"
//	"errors"
//	"fmt"
//	"log"
//	"time"
//)
//
//type Order struct {
//	ID           string
//	User         string
//	Products     []*Product
//	TotalAmount  float64
//	DeliveryAddr string
//	Status       string
//	CreatedAt    time.Time
//}
//
//type OrderModel struct {
//	DB       *sql.DB
//	InfoLog  *log.Logger
//	ErrorLog *log.Logger
//}
//
//func NewOrder(user string, products []*Product, deliveryAddr string) *Order {
//	total := calculateTotal(products)
//	return &Order{
//		User:         user,
//		Products:     products,
//		TotalAmount:  total,
//		DeliveryAddr: deliveryAddr,
//		Status:       "Pending",
//	}
//}
//
//func calculateTotal(products []*Product) float64 {
//	total := 0.0
//	for _, p := range products {
//		total += p.Price
//	}
//	return total
//}
//
//func (o *Order) DisplayOrderDetails() {
//	fmt.Printf("User: %s\n", o.User)
//	fmt.Println("Products:")
//	for _, p := range o.Products {
//		fmt.Printf("- %s (%s) - $%.2f\n", p.Title, p.Description, p.Price)
//	}
//	fmt.Printf("Total amount: $%.2f\n", o.TotalAmount)
//	fmt.Printf("Delivery Address: %s\n", o.DeliveryAddr)
//	fmt.Printf("Status: %s\n", o.Status)
//}
//
//func (m *OrderModel) AddOrder(order *Order) error {
//	// Perform the database insertion
//	_, err := m.DB.Exec("INSERT INTO orders (user, totalamount, deliveryaddr, status, createdat) VALUES ($1, $2, $3, $4, $5)",
//		order.User, order.TotalAmount, order.DeliveryAddr, order.Status, time.Now())
//	if err != nil {
//		m.ErrorLog.Println("Error adding order:", err)
//		return err
//	}
//
//	m.InfoLog.Println("Order added successfully")
//	return nil
//}
//
//func (m *OrderModel) GetOrderByID(id string) (*Order, error) {
//	var order Order
//	err := m.DB.QueryRow("SELECT id, user, totalamount, deliveryaddr, status, createdat FROM orders WHERE id = $1", id).
//		Scan(&order.ID, &order.User, &order.TotalAmount, &order.DeliveryAddr, &order.Status, &order.CreatedAt)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return nil, errors.New("Order not found")
//		}
//		m.ErrorLog.Println("Error getting order:", err)
//		return nil, err
//	}
//	return &order, nil
//}
//
//func (m *OrderModel) GetAllOrders() ([]*Order, error) {
//	rows, err := m.DB.Query("SELECT id, user, totalamount, deliveryaddr, status, createdat FROM orders")
//	if err != nil {
//		m.ErrorLog.Println("Error getting orders:", err)
//		return nil, err
//	}
//	defer rows.Close()
//
//	var orders []*Order
//	for rows.Next() {
//		var order Order
//		if err := rows.Scan(&order.ID, &order.User, &order.TotalAmount, &order.DeliveryAddr, &order.Status, &order.CreatedAt); err != nil {
//			m.ErrorLog.Println("Error scanning order:", err)
//			return nil, err
//		}
//		orders = append(orders, &order)
//	}
//	if err := rows.Err(); err != nil {
//		m.ErrorLog.Println("Error iterating rows:", err)
//		return nil, err
//	}
//
//	return orders, nil
//}
