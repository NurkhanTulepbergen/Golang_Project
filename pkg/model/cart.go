package model

import (
	"database/sql"
	_ "errors"
	"fmt"
	"log"
)

// Struct to represent a shopping cart
type Cart struct {
	UserID string
	Items  map[string]int // Используем map для хранения ID товара и его количества
}

// Function to add a product to the cart
func (c *Cart) AddProduct(productID string, quantity int) {
	// Если товар уже есть в корзине, увеличиваем его количество
	if _, exists := c.Items[productID]; exists {
		c.Items[productID] += quantity
	} else {
		// Иначе добавляем новый товар в корзину
		c.Items[productID] = quantity
	}
}

// Function to remove a product from the cart
func (c *Cart) RemoveProduct(productID string) {
	delete(c.Items, productID)
}

// Function to calculate the total price of items in the cart
func (c *Cart) CalculateTotal(productMap map[string]*Product) float64 {
	total := 0.0
	for productID, quantity := range c.Items {
		if product, ok := productMap[productID]; ok {
			total += float64(quantity) * product.Price
		}
	}
	return total
}

// Function to display the items in the cart
func (c *Cart) DisplayCart(productMap map[string]*Product) {
	fmt.Println("Items in Cart:")
	for productID, quantity := range c.Items {
		if product, ok := productMap[productID]; ok {
			fmt.Printf("- %s (%s) - $%.2f (Quantity: %d)\n", product.Title, product.Description, product.Price, quantity)
		}
	}
	fmt.Printf("Total: $%.2f\n", c.CalculateTotal(productMap))
}

// Struct to represent a shopping cart model
type CartModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

// Method to add a product to the cart in the database
func (m *CartModel) AddProductToCart(userID, productID string, quantity int) error {
	// Perform the database insertion
	_, err := m.DB.Exec("INSERT INTO cart (user_id, product_id, quantity) VALUES ($1, $2, $3)",
		userID, productID, quantity)
	if err != nil {
		m.ErrorLog.Println("Error adding product to cart:", err)
		return err
	}

	m.InfoLog.Println("Product added to cart successfully")
	return nil
}

// Method to remove a product from the cart in the database
func (m *CartModel) RemoveProductFromCart(userID, productID string) error {
	// Perform the database deletion
	_, err := m.DB.Exec("DELETE FROM cart WHERE user_id = $1 AND product_id = $2", userID, productID)
	if err != nil {
		m.ErrorLog.Println("Error removing product from cart:", err)
		return err
	}

	m.InfoLog.Println("Product removed from cart successfully")
	return nil
}

// Method to get the cart for a user from the database
func (m *CartModel) GetCart(userID string) (*Cart, error) {
	rows, err := m.DB.Query("SELECT product_id, quantity FROM cart WHERE user_id = $1", userID)
	if err != nil {
		m.ErrorLog.Println("Error getting cart:", err)
		return nil, err
	}
	defer rows.Close()

	cart := &Cart{UserID: userID, Items: make(map[string]int)}
	for rows.Next() {
		var productID string
		var quantity int
		if err := rows.Scan(&productID, &quantity); err != nil {
			m.ErrorLog.Println("Error scanning cart item:", err)
			return nil, err
		}
		cart.Items[productID] = quantity
	}
	if err := rows.Err(); err != nil {
		m.ErrorLog.Println("Error iterating rows:", err)
		return nil, err
	}

	return cart, nil
}
