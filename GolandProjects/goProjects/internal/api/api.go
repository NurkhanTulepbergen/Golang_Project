package api

import (
	"awesomeProject/pkg/infobook"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Response struct {
	Books []infobook.Info `json:"books"`
}

func StartServer() {
	//create a new router
	router := mux.NewRouter()
	log.Println("creating routes")
	//specify endpoints
	router.HandleFunc("/health-check", HealthCheck).Methods("GET")
	router.HandleFunc("/books", Books).Methods("GET")
	router.HandleFunc("/books/{number}", BookInfo).Methods("GET")
	http.Handle("/", router)

	//start and listen to requests
	http.ListenAndServe(":8080", router)

}
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("entering health check end point")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hey! I am Zhasmin and my application is list of favourite books.")
}

func Books(w http.ResponseWriter, r *http.Request) {
	log.Println("entering books end point")
	var response Response
	books := infobook.PrepareResponse()

	response.Books = books

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return
	}

	w.Write(jsonResponse)
}
func BookInfo(w http.ResponseWriter, r *http.Request) {
	log.Println("entering book info end point")
	vars := mux.Vars(r)
	bookNumber, err := strconv.Atoi(vars["number"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid book number")
		return
	}

	books := infobook.PrepareResponse()

	var foundPlayer *infobook.Info
	for _, book := range books {
		if book.Id == bookNumber {
			foundPlayer = &book
			break
		}
	}

	if foundPlayer == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Book with number %d not found", bookNumber)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, err := json.Marshal(foundPlayer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshalling book info")
		return
	}

	w.Write(jsonResponse)
}
