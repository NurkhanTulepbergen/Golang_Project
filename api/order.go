package api

import (
	"Golang_Project/pkg/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (api *API) CreateOrder(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var order model.Order
	if err := decoder.Decode(&order); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

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

	// Call the GetOrder method of the OrderModel to retrieve the order information
	order, err := api.OrderModel.GetOrder(orderID)
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
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
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
