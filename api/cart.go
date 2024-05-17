package api

import (
	"Golang_Project/pkg/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
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

	// Parse query parameters for filtering, sorting, and pagination
	queryParams := r.URL.Query()
	itemFilter := queryParams.Get("item")
	page, _ := strconv.Atoi(queryParams.Get("page"))
	pageSize, _ := strconv.Atoi(queryParams.Get("pageSize"))
	sortBy := queryParams.Get("sortBy")
	sortOrder := queryParams.Get("sortOrder")

	// Create Filters object with parsed parameters
	filters := model.Filters{
		Item:     itemFilter,
		Page:     page,
		PageSize: pageSize,
		SortBy:   sortBy,
	}

	// Retrieve cart data with applied filters
	carts, metadata, err := api.CartModel.GetCart(filters)
	if err != nil {
		http.Error(w, "Failed to retrieve cart", http.StatusInternalServerError)
		return
	}

	// Apply sorting
	switch sortOrder {
	case "asc":
		switch sortBy {
		case "userId":
			sort.Slice(carts, func(i, j int) bool {
				return carts[i].UserID < carts[j].UserID
			})
			// Add other sorting options if needed
		}
	case "desc":
		switch sortBy {
		case "userId":
			sort.Slice(carts, func(i, j int) bool {
				return carts[i].UserID > carts[j].UserID
			})
			// Add other sorting options if needed
		}
	}

	// Prepare response including metadata
	response := struct {
		Cart     []model.Cart   `json:"cart"`
		Metadata model.Metadata `json:"metadata"`
	}{
		Cart:     carts,
		Metadata: metadata,
	}

	// Set response headers and write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

//func (api *API) UpdateCart(w http.ResponseWriter, r *http.Request) {
//	log.Println("UpdateCart endpoint accessed")
//
//	if r.Method != http.MethodPut {
//		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
//		return
//	}
//
//	var requestData struct {
//		UserID    string `json:"user_id"`
//		ProductID string `json:"product_id"`
//		Quantity  int    `json:"quantity"`
//	}
//
//	err := json.NewDecoder(r.Body).Decode(&requestData)
//	if err != nil {
//		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
//		return
//	}
//
//	err = api.CartModel.UpdateCart(requestData.UserID, requestData.ProductID, requestData.Quantity)
//	if err != nil {
//		http.Error(w, "Failed to update cart", http.StatusInternalServerError)
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//	fmt.Fprintf(w, "Cart updated successfully")
//}
