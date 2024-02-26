package cmd

import (
	"fmt"
)

type Product struct {
	ID   int
	Name string
}

type Shop struct {
	ID   int
	Name string
}

type Models struct {
	Product Product
	Shop    Shop
}

func NewModels() Models {
	return Models{
		Product: Product{
			ID:   1,
			Name: "Product 1",
		},
		Shop: Shop{
			ID:   1,
			Name: "Shop 1",
		},
	}
}

func main() {
	models := NewModels()

	fmt.Println("Product:", models.Product)
	fmt.Println("Shop:", models.Shop)
}
