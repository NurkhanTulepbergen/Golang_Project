package model

import (
	"database/sql"
	"encoding/json"
	_ "errors"
	"fmt"
	"log"
)

type Cart struct {
	UserID string
	Items  map[string]int
}

type CartModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (c *Cart) AddProduct(productID string, quantity int) {

	if _, exists := c.Items[productID]; exists {
		c.Items[productID] += quantity
	} else {

		c.Items[productID] = quantity
	}
}

func (c *Cart) RemoveProduct(productID string) {
	delete(c.Items, productID)
}

func (c *Cart) CalculateTotal(productMap map[string]*Product) float64 {
	total := 0.0
	for productID, quantity := range c.Items {
		if product, ok := productMap[productID]; ok {
			total += float64(quantity) * product.Price
		}
	}
	return total
}

func (c *Cart) DisplayCart(productMap map[string]*Product) {
	fmt.Println("Items in Cart:")
	for productID, quantity := range c.Items {
		if product, ok := productMap[productID]; ok {
			fmt.Printf("- %s (%s) - $%.2f (Quantity: %d)\n", product.Title, product.Description, product.Price, quantity)
		}
	}
	fmt.Printf("Total: $%.2f\n", c.CalculateTotal(productMap))
}

func (m *CartModel) AddProductToCart(userID, productID string, quantity int) error {
	// Формируем строку JSON для нового товара
	newItem := fmt.Sprintf(`{"%s": %d}`, productID, quantity)

	// Выполняем запрос INSERT, чтобы добавить новый товар в корзину пользователя
	_, err := m.DB.Exec("INSERT INTO cart (userid, items) VALUES ($1, $2) ON CONFLICT (userid) DO UPDATE SET items = cart.items || excluded.items",
		userID, newItem)
	if err != nil {
		m.ErrorLog.Println("Error adding product to cart:", err)
		return err
	}

	m.InfoLog.Println("Product added to cart successfully")
	return nil
}

func (m *CartModel) RemoveProductFromCart(userID, productID string) error {
	// Выполняем запрос DELETE для удаления товара из корзины пользователя
	_, err := m.DB.Exec("UPDATE cart SET items = items - $1 WHERE userid = $2",
		fmt.Sprintf(`{"%s": null}`, productID), userID)
	if err != nil {
		m.ErrorLog.Println("Error removing product from cart:", err)
		return err
	}

	m.InfoLog.Println("Product removed from cart successfully")
	return nil
}

func (m *CartModel) GetCart(userID string) (*Cart, error) {
	rows, err := m.DB.Query("SELECT items FROM cart WHERE userid = $1", userID)
	if err != nil {
		m.ErrorLog.Println("Error getting cart:", err)
		return nil, err
	}
	defer rows.Close()

	cart := &Cart{UserID: userID, Items: make(map[string]int)}
	for rows.Next() {
		var itemsJSON []byte
		if err := rows.Scan(&itemsJSON); err != nil {
			m.ErrorLog.Println("Error scanning cart items:", err)
			return nil, err
		}

		var items map[string]int
		if err := json.Unmarshal(itemsJSON, &items); err != nil {
			m.ErrorLog.Println("Error unmarshalling cart items:", err)
			return nil, err
		}

		for productID, quantity := range items {
			cart.Items[productID] = quantity
		}
	}
	if err := rows.Err(); err != nil {
		m.ErrorLog.Println("Error iterating rows:", err)
		return nil, err
	}

	return cart, nil
}
