package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func (api *API) AddProductToCart(w http.ResponseWriter, r *http.Request) {
	log.Println("AddProductToCart endpoint accessed")

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		UserID    string `json:"user_id"`
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	err = api.CartModel.AddProductToCart(requestData.UserID, requestData.ProductID, requestData.Quantity)
	if err != nil {
		http.Error(w, "Failed to add product to cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Product added to cart successfully")
}

// RemoveProductFromCart обрабатывает запрос на удаление товара из корзины.
func (api *API) RemoveProductFromCart(w http.ResponseWriter, r *http.Request) {
	log.Println("RemoveProductFromCart endpoint accessed")

	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		UserID    string `json:"user_id"`
		ProductID string `json:"product_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	err = api.CartModel.RemoveProductFromCart(requestData.UserID, requestData.ProductID)
	if err != nil {
		http.Error(w, "Failed to remove product from cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Product removed from cart successfully")
}

// GetCart обрабатывает запрос на получение содержимого корзины пользователя.
func (api *API) GetCart(w http.ResponseWriter, r *http.Request) {
	log.Println("GetCart endpoint accessed")

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := mux.Vars(r)["user_id"]
	cart, err := api.CartModel.GetCart(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve cart", http.StatusInternalServerError)
		return
	}

	response := struct {
		UserID string         `json:"user_id"`
		Items  map[string]int `json:"items"`
	}{
		UserID: userID,
		Items:  cart.Items,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
