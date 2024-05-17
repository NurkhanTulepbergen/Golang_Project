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
	router.HandleFunc("/shop", api.requireAuthenticatedUser(api.Shops)).Methods("GET")
	router.HandleFunc("/shop", api.requireAuthenticatedUser(api.AddShops)).Methods("POST")
	router.HandleFunc("/shop/{id}", api.requireAuthenticatedUser(api.DeletionByID)).Methods("DELETE")
	router.HandleFunc("/shop/{id}", api.requireAuthenticatedUser(api.UpdateByID)).Methods("PUT")
	router.HandleFunc("/shop/{id}", api.requireAuthenticatedUser(api.GetByID)).Methods("GET")
	router.HandleFunc("/shop/{shop_id}/product", api.requireAuthenticatedUser(api.GetProductsByShopIDHandler)).Methods("GET")

	// Catalog endpoints
	router.HandleFunc("/product", api.requireAuthenticatedUser(api.Products)).Methods("GET")
	router.HandleFunc("/product", api.requireAuthenticatedUser(api.AddProducts)).Methods("POST")
	router.HandleFunc("/product/{id}", api.requireAuthenticatedUser(api.DeleteProductByID)).Methods("DELETE")
	router.HandleFunc("/product/{id}", api.requireAuthenticatedUser(api.UpdateProductByID)).Methods("PUT")
	router.HandleFunc("/product/{id}", api.requireAuthenticatedUser(api.GetProductByID)).Methods("GET")

	// User endpoints
	router.HandleFunc("/user", api.registerUserHandler).Methods("POST")
	router.HandleFunc("/user/activated", api.activateUserHandler).Methods("PUT")

	// Token endpoint
	router.HandleFunc("/tokens/authentication", api.createAuthenticationTokenHandler).Methods("POST")

	router.HandleFunc("/orders", api.requireAuthenticatedUser(api.CreateOrder)).Methods("POST")
	router.HandleFunc("/orders/{order_id}", api.requireAuthenticatedUser(api.GetOrder)).Methods("GET")
	router.HandleFunc("/user/{user_id}/orders", api.requireAuthenticatedUser(api.GetAllOrders)).Methods("GET")
	router.HandleFunc("/orders/{order_id}", api.requireAuthenticatedUser(api.DeleteOrder)).Methods("DELETE")
	router.HandleFunc("/orders/{order_id}", api.requireAuthenticatedUser(api.UpdateOrder)).Methods("PUT")

	router.HandleFunc("/follow/user/{user_id}", api.requireAuthenticatedUser(api.GetFollowDataByUserID)).Methods("GET")
	router.HandleFunc("/follow", api.requireAuthenticatedUser(api.AddProductToFollowList)).Methods("POST")
	router.HandleFunc("/follow/user/{user_id}/product/{product_id}", api.requireAuthenticatedUser(api.DeleteProductFromFollowList)).Methods("DELETE")
	router.HandleFunc("/follow/product/{product_id}", api.requireAuthenticatedUser(api.UpdateProductFromFollowList)).Methods("PUT")

	router.HandleFunc("/history/{userID}", api.GetHistoryHandler).Methods("GET")
	router.HandleFunc("/history", api.AddHistoryHandler).Methods("POST")
	router.HandleFunc("/history/{userID}", api.DeleteHistoryHandler).Methods("DELETE")
	router.HandleFunc("/history", api.UpdateHistoryHandler).Methods("PUT")

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
