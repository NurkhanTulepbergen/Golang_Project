package api

import (
	"Golang_Project/pkg/model"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Response struct {
	Catalog []model.Product `json:"catalog"`
}

func StartServer() {
	router := mux.NewRouter()
	log.Println("creating routes")
	router.HandleFunc("/health-check", HealthCheck).Methods("GET")
	router.HandleFunc("/catalog", Catalog).Methods("GET")
	http.Handle("/", router)
	http.ListenAndServe(":2003", router)
}
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("welcome")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello there")
}
func Catalog(w http.ResponseWriter, r *http.Request) {
	log.Println("books")
	fmt.Fprintf(w, "I love J.S")
}
