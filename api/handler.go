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

type Response struct {
	Shops    []model.Shop    `json:"shops"`
	Products []model.Product `json:"products"`
}

type API struct {
	ShopModel    *model.ShopModel
	ProductModel *model.ProductModel
}

func NewAPI(shopModel *model.ShopModel, productModel *model.ProductModel) *API {
	return &API{ShopModel: shopModel, ProductModel: productModel}
}

func (api *API) StartServer() {
	router := mux.NewRouter()
	log.Println("creating routes")
	router.HandleFunc("/health-check", api.HealthCheck).Methods("GET")
	router.HandleFunc("/shop", api.Shops).Methods("GET")
	router.HandleFunc("/shop", api.AddShops).Methods("POST")
	router.HandleFunc("/shop/{id}", api.DeletionByID).Methods("DELETE")
	router.HandleFunc("/shop/{id}", api.UpdateByID).Methods("PUT")
	router.HandleFunc("/shop/{id}", api.GetByID).Methods("GET")

	router.HandleFunc("/catalog", api.Products).Methods("GET")
	router.HandleFunc("/catalog", api.AddProducts).Methods("POST")
	router.HandleFunc("/catalog/{id}", api.DeleteProductByID).Methods("DELETE")
	router.HandleFunc("/catalog/{id}", api.UpdateProductByID).Methods("PUT")
	router.HandleFunc("/catalog/{id}", api.GetProductByID).Methods("GET")
	http.Handle("/", router)
	http.ListenAndServe(":2003", router)
}

func (api *API) HealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("welcome")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello there")
}

//func (api *API) Shops(w http.ResponseWriter, r *http.Request) {
//	log.Println("getAllShops endpoint accessed")
//
//	if r.Method != http.MethodGet {
//		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
//		return
//	}
//
//	// Получение всех магазинов
//	shops, err := api.ShopModel.GetAllShops()
//	if err != nil {
//		http.Error(w, "Failed to retrieve shops", http.StatusInternalServerError)
//		return
//	}
//	// Формирование ответа в формате JSON
//	response := Response{
//		Shops: shops,
//	}
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	json.NewEncoder(w).Encode(response)
//}

func (api *API) Shops(w http.ResponseWriter, r *http.Request) {
	log.Println("getAllShops endpoint accessed")

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получение всех магазинов
	shops, err := api.ShopModel.GetAllShops()
	if err != nil {
		http.Error(w, "Failed to retrieve shops", http.StatusInternalServerError)
		return
	}

	// Формирование ответа в формате JSON
	response := struct {
		Shops []model.Shop `json:"shops"`
	}{
		Shops: shops,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (api *API) AddShops(w http.ResponseWriter, r *http.Request) {
	log.Println("addShop endpoint accessed")

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the incoming JSON data into a model.Shop struct
	var newShop model.Shop
	err := json.NewDecoder(r.Body).Decode(&newShop)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Call the AddShop method of the ShopModel to add the new shop
	err = api.ShopModel.AddShop(newShop)
	if err != nil {
		http.Error(w, "Failed to add shop", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Shop added successfully")
}
func (api *API) DeletionByID(w http.ResponseWriter, r *http.Request) {
	log.Println("deleteShopByID endpoint accessed")

	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the shop ID from the request URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid shop ID", http.StatusBadRequest)
		return
	}

	// Call the DeleteShopByID method of the ShopModel to delete the shop
	err = api.ShopModel.DeleteShopByID(id)
	if err != nil {
		http.Error(w, "Failed to delete shop", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Shop deleted successfully")
}

func (api *API) UpdateByID(w http.ResponseWriter, r *http.Request) {
	log.Println("updateShopByID endpoint accessed")

	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the shop ID from the request URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid shop ID", http.StatusBadRequest)
		return
	}

	// Decode the request body to get the updated shop data
	var updatedShop model.Shop
	err = json.NewDecoder(r.Body).Decode(&updatedShop)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Call the UpdateShopByID method of the ShopModel to update the shop
	err = api.ShopModel.UpdateShopByID(id, updatedShop)
	if err != nil {
		http.Error(w, "Failed to update shop", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Shop updated successfully")
}

func (api *API) GetByID(w http.ResponseWriter, r *http.Request) {
	log.Println("getShopByID endpoint accessed")

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the shop ID from the request URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Call the GetShopByID method of the ShopModel to retrieve the shop information
	shop, err := api.ShopModel.GetShopByID(id)
	if err != nil {
		http.Error(w, "Failed to get shop", http.StatusInternalServerError)
		return
	}

	// Encode the shop information to JSON
	jsonResponse, err := json.Marshal(shop)
	if err != nil {
		http.Error(w, "Failed to encode shop data", http.StatusInternalServerError)
		return
	}

	// Respond with the shop information
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

//	func (api *API) Products(w http.ResponseWriter, r *http.Request) {
//		log.Println("getAllProducts endpoint accessed")
//
//		if r.Method != http.MethodGet {
//			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
//			return
//		}
//
//		// Получение всех магазинов
//		products, err := api.ProductModel.GetAllProduct()
//		if err != nil {
//			http.Error(w, "Failed to retrieve shops", http.StatusInternalServerError)
//			return
//		}
//		// Формирование ответа в формате JSON
//		response := Response{
//			Products: products,
//		}
//		w.Header().Set("Content-Type", "application/json")
//		w.WriteHeader(http.StatusOK)
//		json.NewEncoder(w).Encode(response)
//	}
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
