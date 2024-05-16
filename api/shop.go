package api

import (
	"Golang_Project/pkg/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sort"
	"strconv"
)

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

func (api *API) GetProductsByShopIDHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GetProductsByShopID endpoint accessed")

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	shopID, err := strconv.ParseInt(mux.Vars(r)["shop_id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid shop ID", http.StatusBadRequest)
		return
	}

	products, err := api.ShopModel.GetProductsByShopID(shopID)
	if err != nil {
		http.Error(w, "Failed to retrieve products for shop", http.StatusInternalServerError)
		return
	}

	// Check if there are no products found
	if len(products) == 0 {
		http.Error(w, "No products found for the provided shop ID", http.StatusNotFound)
		return
	}

	// Encode products as JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "Failed to encode products", http.StatusInternalServerError)
		return
	}
}
