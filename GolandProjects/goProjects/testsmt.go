package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Response struct {
	Books []Info `json:"books"`
}

type Info struct {
	Id     int    `json:"id"`
	Author string `json:"author"`
	Book   string `json:"Book"`
}

func main() {
	log.Println("starting API server")
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
	log.Println("entering persons end point")
	var response Response
	books := prepareResponse()

	response.Books = books

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return
	}

	w.Write(jsonResponse)
}

func prepareResponse() []Info {
	var infos []Info

	var book Info
	book.Id = 1
	book.Author = "Charlotte Bronte"
	book.Book = "Jane Eyr"
	infos = append(infos, book)

	book.Id = 2
	book.Author = "Mikhail Bulgakov"
	book.Book = "Master and Margarita"
	infos = append(infos, book)

	book.Id = 3
	book.Author = "Jack London"
	book.Book = "Martin Iden"
	infos = append(infos, book)

	book.Id = 4
	book.Author = "Ray Bradbury"
	book.Book = "Fahrenheit 451"
	infos = append(infos, book)
	return infos
}
func BookInfo(w http.ResponseWriter, r *http.Request) {
	log.Println("entering player info end point")
	vars := mux.Vars(r)
	bookNumber, err := strconv.Atoi(vars["number"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid book number")
		return
	}

	books := prepareResponse()

	var foundPlayer *Info
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
