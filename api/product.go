package api

import (
	"Golang_Project/pkg/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func (api *API) Products(w http.ResponseWriter, r *http.Request) {
	log.Println("getAllProducts endpoint accessed")

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve all products
	products, err := api.ProductModel.GetAllProduct()
	if err != nil {
		http.Error(w, "Failed to retrieve products", http.StatusInternalServerError)
		return
	}

	// Formulate the response in JSON format
	response := struct {
		Products []model.Product `json:"products"`
	}{
		Products: products,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (api *API) AddProducts(w http.ResponseWriter, r *http.Request) {
	log.Println("addProducts endpoint accessed")

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the incoming JSON data into a model.Shop struct
	var newProduct model.Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Call the AddShop method of the ShopModel to add the new shop
	err = api.ProductModel.AddProduct(newProduct)
	if err != nil {
		http.Error(w, "Failed to add product", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Product added successfully")
}

func (api *API) DeleteProductByID(w http.ResponseWriter, r *http.Request) {
	log.Println("deleteProductByID endpoint accessed")

	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the shop ID from the request URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Call the DeleteShopByID method of the ShopModel to delete the shop
	err = api.ProductModel.DeleteProductByID(id)
	if err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Product deleted successfully")
}

func (api *API) UpdateProductByID(w http.ResponseWriter, r *http.Request) {
	log.Println("updateProductByID endpoint accessed")

	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the shop ID from the request URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Decode the request body to get the updated shop data
	var updatedProduct model.Product
	err = json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Call the UpdateShopByID method of the ShopModel to update the shop
	err = api.ProductModel.UpdateProductByID(id, updatedProduct)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Product updated successfully")
}

func (api *API) GetProductByID(w http.ResponseWriter, r *http.Request) {
	log.Println("getProductByID endpoint accessed")

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the shop ID from the request URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Call the GetShopByID method of the ShopModel to retrieve the shop information
	product, err := api.ProductModel.GetProductByID(id)
	if err != nil {
		http.Error(w, "Failed to get shop", http.StatusInternalServerError)
		return
	}

	// Encode the shop information to JSON
	jsonResponse, err := json.Marshal(product)
	if err != nil {
		http.Error(w, "Failed to encode shop data", http.StatusInternalServerError)
		return
	}

	// Respond with the shop information
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
