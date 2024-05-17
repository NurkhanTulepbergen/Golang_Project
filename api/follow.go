package api

import (
	"Golang_Project/pkg/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (api *API) AddProductToFollowList(w http.ResponseWriter, r *http.Request) {
	log.Println("addFollow endpoint accessed")

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the incoming JSON data into a model.Shop struct
	var flist model.FollowedList
	err := json.NewDecoder(r.Body).Decode(&flist)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Call the AddShop method of the ShopModel to add the new shop
	err = api.FollowModel.AddProductToFollowList(flist)
	if err != nil {
		http.Error(w, "Failed to follow product", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Product followed successfully")

}

func (api *API) DeleteProductFromFollowList(w http.ResponseWriter, r *http.Request) {
	log.Println("deleteFollowProductByID endpoint accessed")

	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the shop ID from the request URL
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	productId, err := strconv.Atoi(vars["product_id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Call the DeleteShopByID method of the ShopModel to delete the shop
	err = api.FollowModel.DeleteProductFromFollowList(userId, productId)
	if err != nil {
		http.Error(w, "Failed to unfollow product", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Product unfollowed successfully")
}

func (api *API) GetFollowDataByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	followData, err := api.FollowModel.GetFollowDataByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to get follow data", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(followData)
	if err != nil {
		http.Error(w, "Failed to encode follow data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
func (api *API) UpdateProductFromFollowList(w http.ResponseWriter, r *http.Request) {
	log.Println("updateFollowProductByID endpoint accessed")

	var flist model.FollowedList
	err := json.NewDecoder(r.Body).Decode(&flist)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Extract the shop ID from the request URL
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["product_id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Call the DeleteShopByID method of the ShopModel to delete the shop
	err = api.FollowModel.UpdateProductFromFollowList(userId, flist)
	if err != nil {
		http.Error(w, "Failed to unfollow product", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Product unfollowed successfully")
}
