package model

import (
	"database/sql"
	"errors"
	"log"
)

type FollowedList struct {
	UserID             int
	ProductID          int
	ProductName        string
	ProductDescription string
	ProductPrice       float64
}

type Follow struct {
	UserID           int
	UserName         string
	FollowedProducts string
}

type FollowModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

// AddProductToFollowList adds a product to the follow list with basic validation.
func (m *FollowModel) AddProductToFollowList(flist FollowedList) error {
	// Basic validation
	if flist.UserID <= 0 {
		return errors.New("invalid UserID: must be a positive integer")
	}
	if flist.ProductID <= 0 {
		return errors.New("invalid ProductID: must be a positive integer")
	}

	_, err := m.DB.Exec(`
        INSERT INTO follow_list (user_id, product_id)
        VALUES ($1, $2)`, flist.UserID, flist.ProductID)
	if err != nil {
		m.ErrorLog.Println("Error adding product:", err)
		return err
	}

	m.InfoLog.Println("Product added successfully")
	return nil
}

// DeleteProductFromFollowList deletes a product from the follow list with basic validation.
func (m *FollowModel) DeleteProductFromFollowList(userID, productID int) error {
	// Basic validation
	if userID <= 0 {
		return errors.New("invalid UserID: must be a positive integer")
	}
	if productID <= 0 {
		return errors.New("invalid ProductID: must be a positive integer")
	}

	_, err := m.DB.Exec(`
        DELETE FROM follow_list 
        WHERE user_id = $1 AND product_id = $2`, userID, productID)
	if err != nil {
		m.ErrorLog.Println("Error deleting product from follow list:", err)
		return err
	}

	m.InfoLog.Println("Product deleted from follow list successfully")
	return nil
}

// GetFollowDataByUserID retrieves follow data for a specific user with basic validation.
func (m *FollowModel) GetFollowDataByUserID(userID int) ([]Follow, error) {
	// Basic validation
	if userID <= 0 {
		return nil, errors.New("invalid UserID: must be a positive integer")
	}

	rows, err := m.DB.Query(`
        SELECT user_id, user_name, followed_products 
        FROM follow
        WHERE user_id = $1`, userID)
	if err != nil {
		m.ErrorLog.Println("Error fetching follow data:", err)
		return nil, err
	}
	defer rows.Close()

	var followData []Follow
	for rows.Next() {
		var data Follow
		if err := rows.Scan(&data.UserID, &data.UserName, &data.FollowedProducts); err != nil {
			m.ErrorLog.Println("Error scanning follow data:", err)
			return nil, err
		}
		followData = append(followData, data)
	}
	if err := rows.Err(); err != nil {
		m.ErrorLog.Println("Error iterating rows:", err)
		return nil, err
	}

	return followData, nil
}

// UpdateProductFromFollowList updates a product in the follow list with basic validation.
func (m *FollowModel) UpdateProductFromFollowList(productsID int, flist FollowedList) error {
	// Basic validation
	if productsID <= 0 {
		return errors.New("invalid ProductID: must be a positive integer")
	}
	if flist.UserID <= 0 {
		return errors.New("invalid UserID: must be a positive integer")
	}
	if flist.ProductID <= 0 {
		return errors.New("invalid ProductID: must be a positive integer")
	}

	// Start a transaction
	tx, err := m.DB.Begin()
	if err != nil {
		m.ErrorLog.Println("Error starting transaction:", err)
		return err
	}
	defer tx.Rollback()

	// Executing UPDATE query to update the product_id in the follow list
	_, err = tx.Exec(`
        UPDATE follow_list 
        SET product_id = $1
        WHERE user_id = $2 AND product_id = $3`, productsID, flist.UserID, flist.ProductID)
	if err != nil {
		m.ErrorLog.Println("Error updating product ID in follow list:", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		m.ErrorLog.Println("Error committing transaction:", err)
		return err
	}

	m.InfoLog.Println("Product ID updated in follow list successfully")
	return nil
}
