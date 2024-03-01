package api

import (
	"Golang_Project/pkg/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type Response struct {
	Shops []model.Shop `json:"shops"`
}

func StartServer() {
	router := mux.NewRouter()
	log.Println("creating routes")
	router.HandleFunc("/health-check", HealthCheck).Methods("GET")
	router.HandleFunc("/shop", Shops).Methods("GET")
	http.Handle("/", router)
	http.ListenAndServe(":2003", router)
}
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("welcome")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello there")
}
func Shops(w http.ResponseWriter, r *http.Request) {
	log.Println("getAllShops endpoint accessed")

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	//shopModel := &model.ShopModel{
	//	DB:       /* Your database connection */,
	//	InfoLog:  log.New(/* Your info logger output */, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
	//	ErrorLog: log.New(/* Your error logger output */, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	//}
	var shops = []model.Shop{
		{Id: "1", CreatedAt: time.Now().String(), UpdatedAt: time.Now().String(), Title: "Shop 1", Description: "Description 1", Type: "Type 1"},
		{Id: "2", CreatedAt: time.Now().String(), UpdatedAt: time.Now().String(), Title: "Shop 2", Description: "Description 2", Type: "Type 2"},
	}

	// Retrieve all shops from the database
	//shops, err := shops.GetAllShops()
	//if err != nil {
	//	http.Error(w, "Failed to retrieve shops", http.StatusInternalServerError)
	//	return
	//}

	response := Response{
		Shops: shops,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
