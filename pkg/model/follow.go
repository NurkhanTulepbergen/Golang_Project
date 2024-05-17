package model

import (
	"database/sql"
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
	FollowedProducts string // Changed to string type
}

type FollowModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m *FollowModel) AddProductToFollowList(flist FollowedList) error {
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

func (m *FollowModel) DeleteProductFromFollowList(userID, productID int) error {
	// Executing DELETE query to remove a product from the user's follow list
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

func (m *FollowModel) GetFollowDataByUserID(userID int) ([]Follow, error) {
	// Executing SELECT query to retrieve follow data for the specified user from the follow_list table
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
