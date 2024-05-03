package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

//func (api *API) StartServer() http.Handler {
//	router := mux.NewRouter()
//	// Convert the app.notFoundResponse helper to a http.Handler using the http.HandlerFunc()
//	// adapter, and then set it as the custom error handler for 404 Not Found responses.
//	router.NotFoundHandler = http.HandlerFunc(api.notFoundResponse)
//
//	// Convert app.methodNotAllowedResponse helper to a http.Handler and set it as the custom
//	// error handler for 405 Method Not Allowed responses
//	router.MethodNotAllowedHandler = http.HandlerFunc(api.methodNotAllowedResponse)
//
//	// healthcheck
//	router.HandleFunc("/health-check", api.requireActivatedUser(api.HealthCheck)).Methods("GET")
//	router.HandleFunc("/shop", api.requireActivatedUser(api.Shops)).Methods("GET")
//	router.HandleFunc("/shop", api.requireActivatedUser(api.AddShops)).Methods("POST")
//	router.HandleFunc("/shop/{id}", api.requireActivatedUser(api.DeletionByID)).Methods("DELETE")
//	router.HandleFunc("/shop/{id}", api.requireActivatedUser(api.UpdateByID)).Methods("PUT")
//	router.HandleFunc("/shop/{id}", api.requireActivatedUser(api.GetByID)).Methods("GET")
//
//	router.HandleFunc("/catalog", api.requireActivatedUser(api.Products)).Methods("GET")
//	router.HandleFunc("/catalog", api.requireActivatedUser(api.AddProducts)).Methods("POST")
//	router.HandleFunc("/catalog/{id}", api.requireActivatedUser(api.DeleteProductByID)).Methods("DELETE")
//	router.HandleFunc("/catalog/{id}", api.requireActivatedUser(api.UpdateProductByID)).Methods("PUT")
//	router.HandleFunc("/catalog/{id}", api.requireActivatedUser(api.GetProductByID)).Methods("GET")
//	router.HandleFunc("/user", api.registerUserHandler).Methods("POST")
//	router.HandleFunc("/user/activated", api.activateUserHandler).Methods("PUT")
//	router.HandleFunc("/tokens/authentication", api.createAuthenticationTokenHandler).Methods("POST")
//	http.Handle("/", api.authenticate(router))
//	http.ListenAndServe(":2003", router)
//	//return api.authenticate(router)
//}

func (api *API) StartServer() {
	router := mux.NewRouter()

	// Health check endpoint
	router.HandleFunc("/health-check", api.requireActivatedUser(api.HealthCheck)).Methods("GET")

	// Shop endpoints
	router.HandleFunc("/shop", api.requirePermission("shop:read", api.Shops)).Methods("GET")
	router.HandleFunc("/shop", api.requirePermission("shop:write", api.AddShops)).Methods("POST")
	router.HandleFunc("/shop/{id}", api.requireActivatedUser(api.DeletionByID)).Methods("DELETE")
	router.HandleFunc("/shop/{id}", api.requireActivatedUser(api.UpdateByID)).Methods("PUT")
	router.HandleFunc("/shop/{id}", api.requireActivatedUser(api.GetByID)).Methods("GET")

	// Catalog endpoints
	router.HandleFunc("/product", api.requireActivatedUser(api.Products)).Methods("GET")
	router.HandleFunc("/product", api.requireActivatedUser(api.AddProducts)).Methods("POST")
	router.HandleFunc("/product/{id}", api.requireActivatedUser(api.DeleteProductByID)).Methods("DELETE")
	router.HandleFunc("/product/{id}", api.requireActivatedUser(api.UpdateProductByID)).Methods("PUT")
	router.HandleFunc("/product/{id}", api.requireActivatedUser(api.GetProductByID)).Methods("GET")

	// User endpoints
	router.HandleFunc("/user", api.registerUserHandler).Methods("POST")
	router.HandleFunc("/user/activated", api.activateUserHandler).Methods("PUT")

	// Token endpoint
	router.HandleFunc("/tokens/authentication", api.createAuthenticationTokenHandler).Methods("POST")

	// Apply middleware
	http.Handle("/", api.authenticate(router))

	// Start the server
	http.ListenAndServe(":2003", nil)
}
