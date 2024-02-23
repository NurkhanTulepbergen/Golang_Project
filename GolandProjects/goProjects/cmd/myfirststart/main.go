package main

import (
	"awesomeProject/internal/api"
	"log"
)

func main() {
	log.Println("starting API server")
	api.StartServer()
}
