package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

// StartServer starts the API server on the specified port.
func (api *API) StartServer(port int) {
	router := mux.NewRouter()

	// Health check endpoint
	router.HandleFunc("/health-check", api.HealthCheck).Methods("GET")

	// Shop endpoints
	router.HandleFunc("/shop", api.Shops).Methods("GET")
	router.HandleFunc("/shop", api.AddShops).Methods("POST")
	router.HandleFunc("/shop/{id}", api.DeletionByID).Methods("DELETE")
	router.HandleFunc("/shop/{id}", api.UpdateByID).Methods("PUT")
	router.HandleFunc("/shop/{id}", api.GetByID).Methods("GET")
	router.HandleFunc("/shop/{shop_id}/product", api.GetProductsByShopIDHandler).Methods("GET")

	// Catalog endpoints
	router.HandleFunc("/product", api.Products).Methods("GET")
	router.HandleFunc("/product", api.AddProducts).Methods("POST")
	router.HandleFunc("/product/{id}", api.DeleteProductByID).Methods("DELETE")
	router.HandleFunc("/product/{id}", api.UpdateProductByID).Methods("PUT")
	router.HandleFunc("/product/{id}", api.GetProductByID).Methods("GET")

	// User endpoints
	router.HandleFunc("/user", api.registerUserHandler).Methods("POST")
	router.HandleFunc("/user/activated", api.activateUserHandler).Methods("PUT")

	// Token endpoint
	router.HandleFunc("/tokens/authentication", api.createAuthenticationTokenHandler).Methods("POST")

	router.HandleFunc("/cart", api.AddProductToCart).Methods("POST")
	router.HandleFunc("/cart", api.RemoveProductFromCart).Methods("PUT")
	router.HandleFunc("/cart", api.GetCart).Methods("GET")

	// Apply middleware
	http.Handle("/", api.authenticate(router))

	// Start the server on the specified port
	addr := fmt.Sprintf(":%d", port)
	http.ListenAndServe(addr, nil)
}

/*// StartServer запускает сервер API на указанном порту.
func (api *API) StartServer() {
	router := mux.NewRouter()

	// Health check endpoint
	router.HandleFunc("/health-check", api.HealthCheck).Methods("GET")

	// Shop endpoints
	router.HandleFunc("/shop", api.Shops).Methods("GET")
	router.HandleFunc("/shop", api.AddShops).Methods("POST")
	router.HandleFunc("/shop/{id}", api.DeletionByID).Methods("DELETE")
	router.HandleFunc("/shop/{id}", api.UpdateByID).Methods("PUT")
	router.HandleFunc("/shop/{id}", api.GetByID).Methods("GET")

	// Catalog endpoints
	router.HandleFunc("/product", api.Products).Methods("GET")
	router.HandleFunc("/product", api.AddProducts).Methods("POST")
	router.HandleFunc("/product/{id}", api.DeleteProductByID).Methods("DELETE")
	router.HandleFunc("/product/{id}", api.UpdateProductByID).Methods("PUT")
	router.HandleFunc("/product/{id}", api.GetProductByID).Methods("GET")

	// User endpoints
	router.HandleFunc("/user", api.registerUserHandler).Methods("POST")
	router.HandleFunc("/user/activated", api.activateUserHandler).Methods("PUT")

	// Token endpoint
	router.HandleFunc("/tokens/authentication", api.createAuthenticationTokenHandler).Methods("POST")
	/*
		router.HandleFunc("/cart", api.AddProductToCart).Methods("POST")
		router.HandleFunc("/cart", api.RemoveProductFromCart).Methods("PUT")
		router.HandleFunc("/cart", api.GetCart).Methods("GET")

	// Запуск маршрутизатора на указанном порту
	go func() {
		if err := http.ListenAndServe(":8080", router); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()
}*/
