package main

import (
	"Golang_Project/api"
	"log"
)

//type Product struct {
//	ID   int
//	Name string
//}
//
//type Shop struct {
//	ID   int
//	Name string
//}
//
//type Models struct {
//	Product Product
//	Shop    Shop
//}
//
//func NewModels() Models {
//	return Models{
//		Product: Product{
//			ID:   1,
//			Name: "Product 1",
//		},
//		Shop: Shop{
//			ID:   1,
//			Name: "Shop 1",
//		},
//	}
//}

func main() {
	log.Println("kickstart my heart")
	api.StartServer()
}
