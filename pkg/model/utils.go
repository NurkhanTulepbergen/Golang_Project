package model

import "sort"

// Filters хранит параметры фильтрации.
type Filters struct {
	Type     string // Пример параметра фильтрации
	SortBy   string // Поле для сортировки
	Page     int
	PageSize int
}

// Metadata хранит метаданные пагинации.
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

// CalculateMetadata вычисляет метаданные пагинации.
func CalculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 || pageSize == 0 {
		return Metadata{} // return an empty Metadata struct if there are no records or pageSize is 0
	}

	lastPage := totalRecords / pageSize
	if totalRecords%pageSize != 0 {
		lastPage++ // increment lastPage if there's a remainder
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     lastPage,
		TotalRecords: totalRecords,
	}
}

// FilterByType фильтрует магазины по типу.
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

// Paginate разбивает магазины на страницы.
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
