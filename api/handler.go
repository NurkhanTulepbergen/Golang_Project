package api

import (
	"Golang_Project/pkg/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Response struct {
	Shops []model.Shop `json:"shops"`
}

type API struct {
	ShopModel *model.ShopModel
}

func NewAPI(shopModel *model.ShopModel) *API {
	return &API{ShopModel: shopModel}
}

func (api *API) StartServer() {
	router := mux.NewRouter()
	log.Println("creating routes")
	router.HandleFunc("/health-check", api.HealthCheck).Methods("GET")
	router.HandleFunc("/shop", api.Shops).Methods("GET")
	http.Handle("/", router)
	http.ListenAndServe(":2003", router)
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

	shops, err := api.ShopModel.GetAllShops()
	if err != nil {
		http.Error(w, "Failed to retrieve shops", http.StatusInternalServerError)
		return
	}

	response := Response{
		Shops: shops,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
