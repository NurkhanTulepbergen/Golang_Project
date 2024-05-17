package api

import (
	"Golang_Project/pkg/model"
	"Golang_Project/pkg/validator"
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

	queryParams := r.URL.Query()
	typeFilter := queryParams.Get("type")
	page, _ := strconv.Atoi(queryParams.Get("page"))
	pageSize, _ := strconv.Atoi(queryParams.Get("pageSize"))
	sortBy := queryParams.Get("sortBy")
	sortOrder := queryParams.Get("sortOrder")

	filters := model.Filters{
		Type:     typeFilter,
		Page:     page,
		PageSize: pageSize,
		SortBy:   sortBy,
	}

	shops, metadata, err := api.ShopModel.GetAllShops(filters)
	if err != nil {
		http.Error(w, "Failed to retrieve shops", http.StatusInternalServerError)
		return
	}

	if sortOrder == "asc" {
		sort.Slice(shops, func(i, j int) bool {
			switch sortBy {
			case "title":
				return shops[i].Title < shops[j].Title
			default:
				return shops[i].Id < shops[j].Id
			}
		})
	} else if sortOrder == "desc" {
		sort.Slice(shops, func(i, j int) bool {
			switch sortBy {
			case "title":
				return shops[i].Title > shops[j].Title
			default:
				return shops[i].Id > shops[j].Id
			}
		})
	}

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

	var newShop model.Shop
	err := json.NewDecoder(r.Body).Decode(&newShop)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	v := validator.New()

	v.Check(newShop.Title != "", "title", "Title is required")
	v.Check(newShop.Description != "", "description", "Description is required")

	if !v.Valid() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(v.Errors)
		return
	}

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

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid shop ID", http.StatusBadRequest)
		return
	}

	err = api.ShopModel.DeleteShopByID(id)
	if err != nil {
		http.Error(w, "Failed to delete shop", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Shop deleted successfully")
}

func (api *API) UpdateByID(w http.ResponseWriter, r *http.Request) {
	log.Println("updateShopByID endpoint accessed")

	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid shop ID", http.StatusBadRequest)
		return
	}

	var updatedShop model.Shop
	err = json.NewDecoder(r.Body).Decode(&updatedShop)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	v := validator.New()

	v.Check(updatedShop.Title != "", "title", "Title is required")
	v.Check(updatedShop.Description != "", "description", "Description is required")

	if !v.Valid() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(v.Errors)
		return
	}

	err = api.ShopModel.UpdateShopByID(id, updatedShop)
	if err != nil {
		http.Error(w, "Failed to update shop", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Shop updated successfully")
}

func (api *API) GetByID(w http.ResponseWriter, r *http.Request) {
	log.Println("getShopByID endpoint accessed")

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	shop, err := api.ShopModel.GetShopByID(id)
	if err != nil {
		http.Error(w, "Failed to get shop", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(shop)
	if err != nil {
		http.Error(w, "Failed to encode shop data", http.StatusInternalServerError)
		return
	}

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

	if len(products) == 0 {
		http.Error(w, "No products found for the provided shop ID", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "Failed to encode products", http.StatusInternalServerError)
		return
	}
}
