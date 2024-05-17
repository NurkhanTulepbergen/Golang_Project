package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"Golang_Project/pkg/model"
	"github.com/gorilla/mux"
	"log"
)

func (api *API) CreateOrder(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var order model.Order
	if err := decoder.Decode(&order); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Calculate the total amount for the order
	total, err := order.CalculateTotal(api.OrderModel.DB)
	if err != nil {
		http.Error(w, "Failed to calculate total amount", http.StatusInternalServerError)
		return
	}

	// Set the total amount for the order
	order.TotalAmount = total

	// Set the creation timestamp
	order.CreatedAt = time.Now()

	if err := api.OrderModel.CreateOrder(&order); err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Order created successfully")
}

func (api *API) GetOrder(w http.ResponseWriter, r *http.Request) {
	log.Println("GetOrder endpoint accessed")

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the order ID from the request URL
	vars := mux.Vars(r)
	orderID := vars["order_id"]

	// Convert orderID to integer
	orderIDInt, err := strconv.Atoi(orderID)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	// Call the GetOrder method of the OrderModel to retrieve the order information
	order, err := api.OrderModel.GetOrder(orderIDInt)
	if err != nil {
		http.Error(w, "Failed to get order", http.StatusInternalServerError)
		return
	}

	// Encode the order information to JSON
	jsonResponse, err := json.Marshal(order)
	if err != nil {
		http.Error(w, "Failed to encode order data", http.StatusInternalServerError)
		return
	}

	// Respond with the order information
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (api *API) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, ok := vars["user_id"]
	if !ok {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Convert userID to integer
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	orders, err := api.OrderModel.GetAllOrders(userIDInt)
	if err != nil {
		http.Error(w, "Failed to get orders", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, "Failed to encode order data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
func (api *API) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderIDStr, ok := vars["order_id"]
	if !ok {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	err = api.OrderModel.DeleteOrder(orderID)
	if err != nil {
		http.Error(w, "Failed to delete order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Order deleted successfully")
}
