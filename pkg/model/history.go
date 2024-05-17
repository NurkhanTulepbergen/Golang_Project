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
type HistoryFilter struct {
	UserID   int    // Фильтр по идентификатору пользователя
	SortBy   string // Поле для сортировки (например, "user_id" или "user_name")
	Order    string // Порядок сортировки ("asc" или "desc")
	Page     int    // Номер страницы
	PageSize int    // Размер страницы
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

func (m *HistoryModel) GetHistory(filter *HistoryFilter) ([]*History, error) {
	// Start building the SQL query
	query := "SELECT user_id, user_name, orders_list FROM history WHERE 1 = 1"
	args := make([]interface{}, 0)

	// Add filters
	if filter.UserID != 0 {
		query += " AND user_id = ?"
		args = append(args, filter.UserID)
	}

	// Add sorting
	if filter.SortBy != "" {
		query += " ORDER BY " + filter.SortBy
		if filter.Order != "" {
			query += " " + filter.Order
		}
	}

	// Add pagination
	if filter.Page != 0 && filter.PageSize != 0 {
		query += " LIMIT ? OFFSET ?"
		args = append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)
	}

	// Execute the query
	rows, err := m.DB.Query(query, args...)
	if err != nil {
		m.ErrorLog.Println("Error getting history:", err)
		return nil, err
	}
	defer rows.Close()

	// Iterate over the result set and build history objects
	historyList := make([]*History, 0)
	for rows.Next() {
		history := &History{}
		var ordersJSON []byte

		err := rows.Scan(&history.UserID, &history.UserName, &ordersJSON)
		if err != nil {
			m.ErrorLog.Println("Error scanning row:", err)
			continue
		}

		err = json.Unmarshal(ordersJSON, &history.OrdersList)
		if err != nil {
			m.ErrorLog.Println("Error unmarshalling orders list:", err)
			continue
		}

		historyList = append(historyList, history)
	}

	return historyList, nil
}

func (m *HistoryModel) UpdateHistory(userID int, userName string, ordersList []*Order) error {
	ordersJSON, err := json.Marshal(ordersList)
	if err != nil {
		m.ErrorLog.Println("Error marshalling orders list:", err)
		return err
	}

	_, err = m.DB.Exec("UPDATE history SET user_name = $1, orders_list = $2 WHERE user_id = $3",
		userName, ordersJSON, userID)
	if err != nil {
		m.ErrorLog.Println("Error updating history:", err)
		return err
	}

	m.InfoLog.Println("History updated successfully")
	return nil
}

func (m *HistoryModel) DeleteHistory(userID int) error {
	_, err := m.DB.Exec("DELETE FROM history WHERE user_id = $1", userID)
	if err != nil {
		m.ErrorLog.Println("Error deleting history:", err)
		return err
	}

	m.InfoLog.Println("History deleted successfully")
	return nil
}
