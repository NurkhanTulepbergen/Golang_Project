package model

import (
	"database/sql"
	"encoding/json"
	"log"
	_ "time"
)

type History struct {
	UserID     int
	UserName   string
	OrdersList []*Order
}

type HistoryModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m *HistoryModel) AddHistory(userID int, userName string, ordersList []*Order) error {
	ordersJSON, err := json.Marshal(ordersList)
	if err != nil {
		m.ErrorLog.Println("Error marshalling orders list:", err)
		return err
	}

	_, err = m.DB.Exec("INSERT INTO history (user_id, user_name, orders_list) VALUES ($1, $2, $3)",
		userID, userName, ordersJSON)
	if err != nil {
		m.ErrorLog.Println("Error adding history:", err)
		return err
	}

	m.InfoLog.Println("History added successfully")
	return nil
}

func (m *HistoryModel) GetHistory(userID int) (*History, error) {
	row := m.DB.QueryRow("SELECT user_id, user_name, orders_list FROM history WHERE user_id = $1", userID)

	history := &History{}
	var ordersJSON []byte

	err := row.Scan(&history.UserID, &history.UserName, &ordersJSON)
	if err != nil {
		m.ErrorLog.Println("Error getting history:", err)
		return nil, err
	}

	err = json.Unmarshal(ordersJSON, &history.OrdersList)
	if err != nil {
		m.ErrorLog.Println("Error unmarshalling orders list:", err)
		return nil, err
	}

	return history, nil
}
