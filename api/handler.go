package api

import (
	"Golang_Project/pkg/model"
	"Golang_Project/pkg/validator"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type Response struct {
	Shops    []model.Shop    `json:"shops"`
	Products []model.Product `json:"products"`
	Users    []model.User    `json:"users"`
	Tokens   []model.Token   `json:"tokens"`
}

type API struct {
	ShopModel    *model.ShopModel
	ProductModel *model.ProductModel
	UserModel    *model.UserModel
	TokenModel   *model.TokenModel
}

func NewAPI(shopModel *model.ShopModel, productModel *model.ProductModel, userModel *model.UserModel, tokenModel *model.TokenModel) *API {
	return &API{ShopModel: shopModel, ProductModel: productModel, UserModel: userModel, TokenModel: tokenModel}
}

func (api *API) StartServer() {
	router := mux.NewRouter()
	log.Println("creating routes")
	router.HandleFunc("/health-check", api.HealthCheck).Methods("GET")
	router.HandleFunc("/shop", api.Shops).Methods("GET")
	router.HandleFunc("/shop", api.AddShops).Methods("POST")
	router.HandleFunc("/shop/{id}", api.DeletionByID).Methods("DELETE")
	router.HandleFunc("/shop/{id}", api.UpdateByID).Methods("PUT")
	router.HandleFunc("/shop/{id}", api.GetByID).Methods("GET")

	router.HandleFunc("/catalog", api.Products).Methods("GET")
	router.HandleFunc("/catalog", api.AddProducts).Methods("POST")
	router.HandleFunc("/catalog/{id}", api.DeleteProductByID).Methods("DELETE")
	router.HandleFunc("/catalog/{id}", api.UpdateProductByID).Methods("PUT")
	router.HandleFunc("/catalog/{id}", api.GetProductByID).Methods("GET")
	router.HandleFunc("/user", api.registerUserHandler).Methods("POST")
	router.HandleFunc("/user/activated", api.activateUserHandler).Methods("PUT")
	router.HandleFunc("/tokens/authentication", api.createAuthenticationTokenHandler).Methods("POST")
	http.Handle("/", router)
	http.ListenAndServe(":2003", router)
}

func (api *API) HealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("welcome")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello there")
}

func (api *API) Shops(w http.ResponseWriter, r *http.Request) {
	log.Println("getAllShops endpoint accessed")

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Разбор параметров запроса для заполнения объекта Filters
	// Для простоты давайте предположим, что параметры запроса используются для фильтрации

	queryParams := r.URL.Query()
	typeFilter := queryParams.Get("type")
	page, _ := strconv.Atoi(queryParams.Get("page"))
	pageSize, _ := strconv.Atoi(queryParams.Get("pageSize"))
	sortBy := queryParams.Get("sortBy")
	sortOrder := queryParams.Get("sortOrder")

	// Создание объекта Filters с разобранными параметрами
	filters := model.Filters{
		Type:     typeFilter,
		Page:     page,
		PageSize: pageSize,
		SortBy:   sortBy,
	}

	// Получение магазинов с примененными фильтрами
	shops, metadata, err := api.ShopModel.GetAllShops(filters)
	if err != nil {
		http.Error(w, "Failed to retrieve shops", http.StatusInternalServerError)
		return
	}

	// Применение сортировки
	if sortOrder == "asc" {
		sort.Slice(shops, func(i, j int) bool {
			switch sortBy {
			case "title":
				return shops[i].Title < shops[j].Title
			// Добавьте другие варианты сортировки при необходимости
			default:
				return shops[i].Id < shops[j].Id
			}
		})
	} else if sortOrder == "desc" {
		sort.Slice(shops, func(i, j int) bool {
			switch sortBy {
			case "title":
				return shops[i].Title > shops[j].Title
			// Добавьте другие варианты сортировки при необходимости
			default:
				return shops[i].Id > shops[j].Id
			}
		})
	}

	// Формирование ответа включая метаданные
	response := struct {
		Shops    []model.Shop   `json:"shops"`
		Metadata model.Metadata `json:"metadata"`
	}{
		Shops:    shops,
		Metadata: metadata,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

//func (api *API) Shops(w http.ResponseWriter, r *http.Request) {
//	log.Println("getAllShops endpoint accessed")
//
//	if r.Method != http.MethodGet {
//		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
//		return
//	}
//
//	// Разбор параметров запроса для заполнения объекта Filters
//	// Для простоты давайте предположим, что параметры запроса используются для фильтрации
//
//	queryParams := r.URL.Query()
//	typeFilter := queryParams.Get("type")
//	page, _ := strconv.Atoi(queryParams.Get("page"))
//	pageSize, _ := strconv.Atoi(queryParams.Get("pageSize"))
//	sortBy := queryParams.Get("sortBy")
//
//	// Создание объекта Filters с разобранными параметрами
//	filters := model.Filters{
//		Type:     typeFilter,
//		Page:     page,
//		PageSize: pageSize,
//		SortBy:   sortBy,
//	}
//
//	// Получение магазинов с примененными фильтрами
//	shops, metadata, err := api.ShopModel.GetAllShops(filters)
//	if err != nil {
//		http.Error(w, "Failed to retrieve shops", http.StatusInternalServerError)
//		return
//	}
//
//	// Формирование ответа включая метаданные
//	response := struct {
//		Shops    []model.Shop   `json:"shops"`
//		Metadata model.Metadata `json:"metadata"`
//	}{
//		Shops:    shops,
//		Metadata: metadata,
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	json.NewEncoder(w).Encode(response)
//}

func (api *API) AddShops(w http.ResponseWriter, r *http.Request) {
	log.Println("addShop endpoint accessed")

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the incoming JSON data into a model.Shop struct
	var newShop model.Shop
	err := json.NewDecoder(r.Body).Decode(&newShop)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Call the AddShop method of the ShopModel to add the new shop
	err = api.ShopModel.AddShop(newShop)
	if err != nil {
		http.Error(w, "Failed to add shop", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Shop added successfully")
}
func (api *API) DeletionByID(w http.ResponseWriter, r *http.Request) {
	log.Println("deleteShopByID endpoint accessed")

	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the shop ID from the request URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid shop ID", http.StatusBadRequest)
		return
	}

	// Call the DeleteShopByID method of the ShopModel to delete the shop
	err = api.ShopModel.DeleteShopByID(id)
	if err != nil {
		http.Error(w, "Failed to delete shop", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Shop deleted successfully")
}

func (api *API) UpdateByID(w http.ResponseWriter, r *http.Request) {
	log.Println("updateShopByID endpoint accessed")

	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the shop ID from the request URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid shop ID", http.StatusBadRequest)
		return
	}

	// Decode the request body to get the updated shop data
	var updatedShop model.Shop
	err = json.NewDecoder(r.Body).Decode(&updatedShop)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Call the UpdateShopByID method of the ShopModel to update the shop
	err = api.ShopModel.UpdateShopByID(id, updatedShop)
	if err != nil {
		http.Error(w, "Failed to update shop", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Shop updated successfully")
}

func (api *API) GetByID(w http.ResponseWriter, r *http.Request) {
	log.Println("getShopByID endpoint accessed")

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the shop ID from the request URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Call the GetShopByID method of the ShopModel to retrieve the shop information
	shop, err := api.ShopModel.GetShopByID(id)
	if err != nil {
		http.Error(w, "Failed to get shop", http.StatusInternalServerError)
		return
	}

	// Encode the shop information to JSON
	jsonResponse, err := json.Marshal(shop)
	if err != nil {
		http.Error(w, "Failed to encode shop data", http.StatusInternalServerError)
		return
	}

	// Respond with the shop information
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

//	func (api *API) Products(w http.ResponseWriter, r *http.Request) {
//		log.Println("getAllProducts endpoint accessed")
//
//		if r.Method != http.MethodGet {
//			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
//			return
//		}
//
//		// Получение всех магазинов
//		products, err := api.ProductModel.GetAllProduct()
//		if err != nil {
//			http.Error(w, "Failed to retrieve shops", http.StatusInternalServerError)
//			return
//		}
//		// Формирование ответа в формате JSON
//		response := Response{
//			Products: products,
//		}
//		w.Header().Set("Content-Type", "application/json")
//		w.WriteHeader(http.StatusOK)
//		json.NewEncoder(w).Encode(response)
//	}
func (api *API) Products(w http.ResponseWriter, r *http.Request) {
	log.Println("getAllProducts endpoint accessed")

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve all products
	products, err := api.ProductModel.GetAllProduct()
	if err != nil {
		http.Error(w, "Failed to retrieve products", http.StatusInternalServerError)
		return
	}

	// Formulate the response in JSON format
	response := struct {
		Products []model.Product `json:"products"`
	}{
		Products: products,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (api *API) AddProducts(w http.ResponseWriter, r *http.Request) {
	log.Println("addProducts endpoint accessed")

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the incoming JSON data into a model.Shop struct
	var newProduct model.Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Call the AddShop method of the ShopModel to add the new shop
	err = api.ProductModel.AddProduct(newProduct)
	if err != nil {
		http.Error(w, "Failed to add product", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Product added successfully")
}

func (api *API) DeleteProductByID(w http.ResponseWriter, r *http.Request) {
	log.Println("deleteProductByID endpoint accessed")

	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the shop ID from the request URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Call the DeleteShopByID method of the ShopModel to delete the shop
	err = api.ProductModel.DeleteProductByID(id)
	if err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Product deleted successfully")
}

func (api *API) UpdateProductByID(w http.ResponseWriter, r *http.Request) {
	log.Println("updateProductByID endpoint accessed")

	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the shop ID from the request URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Decode the request body to get the updated shop data
	var updatedProduct model.Product
	err = json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Call the UpdateShopByID method of the ShopModel to update the shop
	err = api.ProductModel.UpdateProductByID(id, updatedProduct)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Product updated successfully")
}

func (api *API) GetProductByID(w http.ResponseWriter, r *http.Request) {
	log.Println("getProductByID endpoint accessed")

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the shop ID from the request URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Call the GetShopByID method of the ShopModel to retrieve the shop information
	product, err := api.ProductModel.GetProductByID(id)
	if err != nil {
		http.Error(w, "Failed to get shop", http.StatusInternalServerError)
		return
	}

	// Encode the shop information to JSON
	jsonResponse, err := json.Marshal(product)
	if err != nil {
		http.Error(w, "Failed to encode shop data", http.StatusInternalServerError)
		return
	}

	// Respond with the shop information
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
func (api *API) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the incoming JSON data into a struct
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Create a new User instance
	newuser := &model.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}
	// Use the Password.Set() method to generate and store the hashed and plaintext
	// passwords.
	err = newuser.Password.Set(input.Password)
	if err != nil {
		http.Error(w, "Failed to set password", http.StatusInternalServerError)
		return
	}

	v := validator.New()
	// Validate the user struct
	model.ValidateUser(v, newuser)
	if !v.Valid() {
		errMsg, _ := json.Marshal(v.Errors)
		http.Error(w, string(errMsg), http.StatusBadRequest)
		return
	}

	// Insert the user data into the database
	err = api.UserModel.Insert(newuser)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrDuplicateEmail):
			http.Error(w, "User with this email already exists", http.StatusBadRequest)
			return
		default:
			http.Error(w, "Failed to insert user: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// Respond with success message
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User registered successfully")

	if api.TokenModel == nil {
		http.Error(w, "Token model is not initialized", http.StatusInternalServerError)
		return
	}

	token, err := api.TokenModel.New(newuser.ID, 3*24*time.Hour, "activated")
	if err != nil {
		http.Error(w, "Failed to generate activation token", http.StatusInternalServerError)
		return
	}

	// Create the response struct
	response := struct {
		User  *model.User `json:"user"`
		Token string      `json:"token"`
	}{
		User:  newuser,
		Token: token.Plaintext,
	}

	// Encode the response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	// Write the response to the client
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}

}

//	func (api *API) activateUserHandler(w http.ResponseWriter, r *http.Request) {
//		// Parse the plaintext activation token from the request body.
//		var input struct {
//			TokenPlaintext string `json:"token"`
//		}
//		err := json.NewDecoder(r.Body).Decode(&input)
//		if err != nil {
//			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
//			return
//		}
//
//		// Validate the plaintext token provided by the client.
//		v := validator.New()
//		if model.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
//			errMsg, _ := json.Marshal(v.Errors)
//			http.Error(w, string(errMsg), http.StatusBadRequest)
//			return
//		}
//
//		// Retrieve the details of the user associated with the token using the
//		// GetForToken() method (which we will create in a minute). If no matching record
//		// is found, then we let the client know that the token they provided is not valid.
//		newuser, err := api.UserModel.GetForToken(model.ScopeActivation, input.TokenPlaintext)
//		if err != nil {
//			switch {
//			case errors.Is(err, model.ErrRecordNotFound):
//				// Check if the error is model.ErrRecordNotFound, indicating that no user record was found for the provided token.
//				// Since the error is not specific to token expiration, we should not assume that the token is expired.
//				// Instead, we can simply return a message indicating that the token is invalid.
//				v.AddError("token", "invalid activation token")
//				errMsg, _ := json.Marshal(v.Errors)
//				http.Error(w, string(errMsg), http.StatusBadRequest)
//			default:
//				// For any other error, return a generic error message indicating a problem with retrieving the user for the token.
//				http.Error(w, "Failed to get user for token", http.StatusInternalServerError)
//			}
//			//switch {
//			//case errors.Is(err, model.ErrRecordNotFound):
//			//	v.AddError("token", "invalid or expired activation token")
//			//	errMsg, _ := json.Marshal(v.Errors)
//			//	http.Error(w, string(errMsg), http.StatusBadRequest)
//			//default:
//			//	http.Error(w, "Failed to get user for token", http.StatusInternalServerError)
//			//}
//			return
//		}
//
//		if newuser == nil {
//			// Log a message if newuser is nil
//			log.Println("Error: newuser is nil")
//			http.Error(w, "Failed to activate user: user data not found", http.StatusInternalServerError)
//			return
//		}
//
//		// Update the user's activation status.
//		newuser.Activated = true
//
//		// Save the updated user record in our database, checking for any edit conflicts in
//		// the same way that we did for our movie records.
//		err = api.UserModel.Update(newuser)
//		if err != nil {
//			http.Error(w, "Failed to update user", http.StatusInternalServerError)
//			return
//		}
//
//		// If everything went successfully, then we delete all activation tokens for the
//		// user.
//		err = api.TokenModel.DeleteAllForUser(model.ScopeActivation, newuser.ID)
//		if err != nil {
//			http.Error(w, "Failed to delete activation tokens", http.StatusInternalServerError)
//			return
//		}
//
//		// Send the updated user details to the client in a JSON response.
//		response := struct {
//			User *model.User `json:"user"`
//		}{
//			User: newuser,
//		}
//
//		jsonResponse, err := json.Marshal(response)
//		if err != nil {
//			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
//			return
//		}
//
//		w.Header().Set("Content-Type", "application/json")
//		w.WriteHeader(http.StatusOK)
//		_, err = w.Write(jsonResponse)
//		if err != nil {
//			http.Error(w, "Failed to write response", http.StatusInternalServerError)
//			return
//		}
//	}
//
//	func (api *API) activateUserHandler(w http.ResponseWriter, r *http.Request) {
//		// Parse the plaintext activation token from the request body.
//		var input struct {
//			TokenPlaintext string `json:"token"`
//		}
//		err := api.readJSON(w, r, &input)
//		if err != nil {
//			api.badRequestResponse(w, r, err)
//			return
//		}
//		// Validate the plaintext token provided by the client.
//		v := validator.New()
//		if model.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
//			api.failedValidationResponse(w, r, v.Errors)
//			return
//		}
//		// Retrieve the details of the user associated with the token using the
//		// GetForToken() method (which we will create in a minute). If no matching record
//		// is found, then we let the client know that the token they provided is not valid.
//		//user, err := api.UserModel.GetForToken(model.ScopeActivation, input.TokenPlaintext)
//		//if err != nil {
//		//	switch {
//		//	case errors.Is(err, model.ErrRecordNotFound):
//		//		v.AddError("token", "invalid or expired activation token")
//		//		api.failedValidationResponse(w, r, v.Errors)
//		//	//case errors.Is(err, model.ErrTokenExpired):
//		//	//	v.AddError("token", "expired activation token")
//		//	//	api.failedValidationResponse(w, r, v.Errors)
//		//	default:
//		//		api.serverErrorResponse(w, r, err)
//		//	}
//		//	return
//		//}
//		// Retrieve the user associated with the token.
//		user, err := api.UserModel.GetForToken(model.ScopeActivation, input.TokenPlaintext)
//
//		// Check if the user was found.
//		if user == nil {
//			// If no user is found, then we let the client know that the token they provided is not valid.
//			v := validator.New()
//			v.AddError("token", "invalid or expired activation token")
//			api.failedValidationResponse(w, r, v.Errors)
//			return
//		}
//		if err != nil {
//			if errors.Is(err, model.ErrRecordNotFound) {
//				// If no matching record is found, then we let the client know that the token they provided is not valid.
//				v := validator.New()
//				v.AddError("token", "invalid or expired activation token")
//				api.failedValidationResponse(w, r, v.Errors)
//				return
//			}
//			api.serverErrorResponse(w, r, err)
//			return
//		}
//		// Update the user's activation status.
//		user.Activated = true
//		// Save the updated user record in our database, checking for any edit conflicts in
//		// the same way that we did for our movie records.
//		err = api.UserModel.Update(user)
//		if err != nil {
//			api.serverErrorResponse(w, r, err)
//			return
//		}
//		// If everything went successfully, then we delete all activation tokens for the
//		// user.
//		err = api.TokenModel.DeleteAllForUser(model.ScopeActivation, user.ID)
//		if err != nil {
//			api.serverErrorResponse(w, r, err)
//			return
//		}
//		// Send the updated user details to the client in a JSON response.
//		err = api.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
//		if err != nil {
//			api.serverErrorResponse(w, r, err)
//		}
//
// }
func (api *API) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the plaintext activation token from the request body.
	var input struct {
		TokenPlaintext string `json:"token"`
	}
	err := api.readJSON(w, r, &input)
	if err != nil {
		api.badRequestResponse(w, r, err)
		return
	}
	// Validate the plaintext token provided by the client.
	v := validator.New()
	if model.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		api.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Retrieve the details of the user associated with the token using the
	// GetForToken() method (which we will create in a minute). If no matching record
	// is found, then we let the client know that the token they provided is not valid.
	user, err := api.UserModel.GetForToken(model.ScopeActivation, input.TokenPlaintext)

	// Check if the user was found.
	if user == nil {
		// If no user is found, then we let the client know that the token they provided is not valid.
		v := validator.New()
		v.AddError("token", "invalid or expired activation token")
		api.failedValidationResponse(w, r, v.Errors)
		return
	}
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			// If no matching record is found, then we let the client know that the token they provided is not valid.
			v := validator.New()
			v.AddError("token", "invalid or expired activation token")
			api.failedValidationResponse(w, r, v.Errors)
			return
		}
		api.serverErrorResponse(w, r, err)
		return
	}
	// Update the user's activation status.
	user.Activated = true
	// Save the updated user record in our database, checking for any edit conflicts in
	// the same way that we did for our movie records.
	err = api.UserModel.Update(user)
	if err != nil {
		api.serverErrorResponse(w, r, err)
		return
	}
	// If everything went successfully, then we delete all activation tokens for the
	// user.
	err = api.TokenModel.DeleteAllForUser(model.ScopeActivation, user.ID)
	if err != nil {
		api.serverErrorResponse(w, r, err)
		return
	}
	// Send the updated user details to the client in a JSON response.
	err = api.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		api.serverErrorResponse(w, r, err)
		return
	}
	// Return success without error
	return
}

func (api *API) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the email and password from the request body.
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := api.readJSON(w, r, &input)
	if err != nil {
		api.badRequestResponse(w, r, err)
		return
	}
	// Validate the email and password provided by the client.
	v := validator.New()
	model.ValidateEmail(v, input.Email)
	model.ValidatePasswordPlaintext(v, input.Password)
	if !v.Valid() {
		api.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Lookup the user record based on the email address. If no matching user was
	// found, then we call the app.invalidCredentialsResponse() helper to send a 401
	// Unauthorized response to the client (we will create this helper in a moment).
	user, err := api.UserModel.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			api.invalidCredentialsResponse(w, r)
		default:
			api.serverErrorResponse(w, r, err)
		}
		return
	}
	// Check if the provided password matches the actual password for the user.
	match, err := user.Password.Matches(input.Password)
	if err != nil {
		api.serverErrorResponse(w, r, err)
		return
	}
	// If the passwords don't match, then we call the app.invalidCredentialsResponse()
	// helper again and return.
	if !match {
		api.invalidCredentialsResponse(w, r)
		return
	}
	// Otherwise, if the password is correct, we generate a new token with a 24-hour
	// expiry time and the scope 'authentication'.
	token, err := api.TokenModel.New(user.ID, 24*time.Hour, model.ScopeAuthentication)
	if err != nil {
		api.serverErrorResponse(w, r, err)
		return
	}
	// Encode the token to JSON and send it in the response along with a 201 Created
	// status code.
	err = api.writeJSON(w, http.StatusCreated, envelope{
		"authentication_token": token,
	}, nil)
	if err != nil {
		api.serverErrorResponse(w, r, err)
	}
}
