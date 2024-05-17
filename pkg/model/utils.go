package model

import "sort"

type Filters struct {
	Item      string
	Title     string
	Type      string // Пример параметра фильтрации
	SortBy    string // Поле для сортировки
	Page      int
	PageSize  int
	SortOrder string
}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func CalculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 || pageSize == 0 {
		return Metadata{}
	}

	lastPage := totalRecords / pageSize
	if totalRecords%pageSize != 0 {
		lastPage++
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     lastPage,
		TotalRecords: totalRecords,
	}
}

func FilterByType(shops []Shop, shopType string) []Shop {
	var filteredShops []Shop
	for _, shop := range shops {
		if shop.Type == shopType {
			filteredShops = append(filteredShops, shop)
		}
	}
	return filteredShops
}

func SortByTitle(shops []Shop) []Shop {
	sort.Slice(shops, func(i, j int) bool {
		return shops[i].Title < shops[j].Title
	})
	return shops
}

func Paginate(shops []Shop, page, pageSize int) []Shop {
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(shops) {
		return nil
	}

	if end > len(shops) {
		end = len(shops)
	}

	return shops[start:end]
}
func FilterByTitle(products []Product, productTitle string) []Product {
	var filteredProducts []Product
	for _, product := range products {
		if product.Title == productTitle {
			filteredProducts = append(filteredProducts, product)
		}
	}
	return filteredProducts
}

func SortByPrice(products []Product, sortBy string) []Product {
	switch sortBy {
	case "price":
		sort.Slice(products, func(i, j int) bool {
			return products[i].Price < products[j].Price
		})
		// Add more cases for additional fields if needed
	}
	return products
}

func PaginateForProduct(product []Product, page, pageSize int) []Product {
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(product) {
		return nil
	}

	if end > len(product) {
		end = len(product)
	}

	return product[start:end]
}

func FilterByItems(carts []Cart, cartItem string) []Cart {
	var filteredCart []Cart
	for _, cart := range carts {
		// Check if the cart contains the specified item
		if _, ok := cart.Items[cartItem]; ok {
			filteredCart = append(filteredCart, cart)
		}
	}
	return filteredCart
}

func SortById(carts []Cart, sortBy string) []Cart {
	switch sortBy {
	case "userId":
		sort.Slice(carts, func(i, j int) bool {
			return carts[i].UserID < carts[j].UserID
		})
		// Add more cases for additional fields if needed
	}
	return carts
}

func PaginateForCarts(carts []Cart, page, pageSize int) []Cart {
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(carts) {
		return nil
	}

	if end > len(carts) {
		end = len(carts)
	}

	return carts[start:end]
}
