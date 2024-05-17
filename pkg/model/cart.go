package model

import (
	"database/sql"
	"errors"
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

func (m *CartModel) GetCart(filters Filters) ([]Cart, Metadata, error) {
	// Construct the base SQL query
	query := "SELECT items, userid FROM cart"

	// Execute the SQL query
	rows, err := m.DB.Query(query)
	if err != nil {
		m.ErrorLog.Println("Error getting cart:", err)
		return nil, Metadata{}, err
	}
	defer rows.Close()

	var carts []Cart
	for rows.Next() {
		var cart Cart
		if err := rows.Scan(&cart.Items, &cart.UserID); err != nil {
			m.ErrorLog.Println("Error scanning cart:", err)
			return nil, Metadata{}, err
		}
		carts = append(carts, cart)
	}
	if err := rows.Err(); err != nil {
		m.ErrorLog.Println("Error iterating rows:", err)
		return nil, Metadata{}, err
	}

	// Apply filtering if necessary
	if filters.Item != "" {
		carts = FilterByItems(carts, filters.Item)
	}

	// Apply sorting if necessary
	if filters.SortBy != "" {
		switch filters.SortBy {
		case "userId":
			carts = SortById(carts, filters.SortBy)
		default:
			return nil, Metadata{}, errors.New("unknown sort field")
		}
	}

	// Paginate the cart slice
	paginatedCart := PaginateForCarts(carts, filters.Page, filters.PageSize)

	// Calculate metadata based on the number of records after filtering and pagination
	totalRecords := len(carts)
	metadata := CalculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return paginatedCart, metadata, nil
}

//func (m *CartModel) UpdateCart(userID, productID string, newQuantity int) error {
//	// Fetch the current cart of the user
//	cart, err, _ := m.GetCart(userID)
//	if err != nil {
//		return err
//	}
//
//	// Check if the product already exists in the cart
//	_, exists := cart.Items[productID]
//
//	// If the product doesn't exist, add it to the cart with the new quantity
//	if !exists {
//		cart.AddProduct(productID, newQuantity)
//	} else {
//		// Update the quantity of the existing product in the cart
//		if newQuantity <= 0 {
//			// If the new quantity is zero or negative, remove the product from the cart
//			delete(cart.Items, productID)
//		} else {
//			// Otherwise, update the quantity
//			cart.Items[productID] = newQuantity
//		}
//	}
//
//	// Marshal the cart items into JSON
//	itemsJSON, err := json.Marshal(cart.Items)
//	if err != nil {
//		m.ErrorLog.Println("Error marshalling cart items:", err)
//		return err
//	}
//
//	// Update the cart in the database
//	_, err = m.DB.Exec("UPDATE cart SET items = $1 WHERE userid = $2", itemsJSON, userID)
//	if err != nil {
//		m.ErrorLog.Println("Error updating cart:", err)
//		return err
//	}
//
//	m.InfoLog.Println("Cart updated successfully")
//	return nil
//}
