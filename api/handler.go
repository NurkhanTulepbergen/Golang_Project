package api

import (
	"Golang_Project/pkg/model"
	//"Golang_Project/pkg/validator"
	//"context"
	"encoding/json"
	//"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sort"
	"strconv"
	//"time"
)

type Response struct {
	Shops       []model.Shop        `json:"shops"`
	Products    []model.Product     `json:"products"`
	Users       []model.User        `json:"users"`
	Tokens      []model.Token       `json:"tokens"`
	Permissions []model.Permissions `json:"permissions"`
	Carts       []model.Cart        `json:"carts"`
}

type API struct {
	ShopModel       *model.ShopModel
	ProductModel    *model.ProductModel
	UserModel       *model.UserModel
	TokenModel      *model.TokenModel
	PermissionModel *model.PermissionModel
	CartModel       *model.CartModel
}

func NewAPI(shopModel *model.ShopModel, productModel *model.ProductModel, userModel *model.UserModel, tokenModel *model.TokenModel, permissionModel *model.PermissionModel, cartModel *model.CartModel) *API {
	return &API{ShopModel: shopModel, ProductModel: productModel, UserModel: userModel, TokenModel: tokenModel, PermissionModel: permissionModel, CartModel: cartModel}
}

func (api *API) HealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("welcome")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello there")
}

func (api *API) Shops(w http.ResponseWriter, r *http.Request) {
	log.Println("getAllShops endpoint accessed")

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Разбор параметров запроса для заполнения объекта Filters
	// Для простоты давайте предположим, что параметры запроса используются для фильтрации

	queryParams := r.URL.Query()
	typeFilter := queryParams.Get("type")
	page, _ := strconv.Atoi(queryParams.Get("page"))
	pageSize, _ := strconv.Atoi(queryParams.Get("pageSize"))
	sortBy := queryParams.Get("sortBy")
	sortOrder := queryParams.Get("sortOrder")

	// Создание объекта Filters с разобранными параметрами
	filters := model.Filters{
		Type:     typeFilter,
		Page:     page,
		PageSize: pageSize,
		SortBy:   sortBy,
	}

	// Получение магазинов с примененными фильтрами
	shops, metadata, err := api.ShopModel.GetAllShops(filters)
	if err != nil {
		http.Error(w, "Failed to retrieve shops", http.StatusInternalServerError)
		return
	}

	// Применение сортировки
	if sortOrder == "asc" {
		sort.Slice(shops, func(i, j int) bool {
			switch sortBy {
			case "title":
				return shops[i].Title < shops[j].Title
			// Добавьте другие варианты сортировки при необходимости
			default:
				return shops[i].Id < shops[j].Id
			}
		})
	} else if sortOrder == "desc" {
		sort.Slice(shops, func(i, j int) bool {
			switch sortBy {
			case "title":
				return shops[i].Title > shops[j].Title
			// Добавьте другие варианты сортировки при необходимости
			default:
				return shops[i].Id > shops[j].Id
			}
		})
	}

	// Формирование ответа включая метаданные
	response := struct {
		Shops    []model.Shop   `json:"shops"`
		Metadata model.Metadata `json:"metadata"`
	}{
		Shops:    shops,
		Metadata: metadata,
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

// AddProductToCart обрабатывает запрос на добавление товара в корзину.
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
