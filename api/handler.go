package api

import (
	"Golang_Project/pkg/model"
	"Golang_Project/pkg/validator"
	"context"
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
	Shops       []model.Shop        `json:"shops"`
	Products    []model.Product     `json:"products"`
	Users       []model.User        `json:"users"`
	Tokens      []model.Token       `json:"tokens"`
	Permissions []model.Permissions `json:"permissions"`
}

type API struct {
	ShopModel       *model.ShopModel
	ProductModel    *model.ProductModel
	UserModel       *model.UserModel
	TokenModel      *model.TokenModel
	PermissionModel *model.PermissionModel
}

func NewAPI(shopModel *model.ShopModel, productModel *model.ProductModel, userModel *model.UserModel, tokenModel *model.TokenModel, permissionModel *model.PermissionModel) *API {
	return &API{ShopModel: shopModel, ProductModel: productModel, UserModel: userModel, TokenModel: tokenModel, PermissionModel: permissionModel}
}

//func (api *API) StartServer() {
//	router := mux.NewRouter()
//	log.Println("creating routes")
//	router.HandleFunc("/health-check", api.HealthCheck).Methods("GET")
//	router.HandleFunc("/shop", api.Shops).Methods("GET")
//	router.HandleFunc("/shop", api.AddShops).Methods("POST")
//	router.HandleFunc("/shop/{id}", api.DeletionByID).Methods("DELETE")
//	router.HandleFunc("/shop/{id}", api.UpdateByID).Methods("PUT")
//	router.HandleFunc("/shop/{id}", api.GetByID).Methods("GET")
//
//	router.HandleFunc("/catalog", api.Products).Methods("GET")
//	router.HandleFunc("/catalog", api.AddProducts).Methods("POST")
//	router.HandleFunc("/catalog/{id}", api.DeleteProductByID).Methods("DELETE")
//	router.HandleFunc("/catalog/{id}", api.UpdateProductByID).Methods("PUT")
//	router.HandleFunc("/catalog/{id}", api.GetProductByID).Methods("GET")
//	router.HandleFunc("/user", api.registerUserHandler).Methods("POST")
//	router.HandleFunc("/user/activated", api.activateUserHandler).Methods("PUT")
//	router.HandleFunc("/tokens/authentication", api.createAuthenticationTokenHandler).Methods("POST")
//	http.Handle("/", router)
//	http.ListenAndServe(":2003", router)
//
//}

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
	// Create an anonymous struct to hold the expected data from the request body.
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// Parse the request body into the anonymous struct.
	err := api.readJSON(w, r, &input)
	if err != nil {
		api.badRequestResponse(w, r, err)
		return
	}
	// Copy the data from the request body into a new User struct. Notice also that we
	// set the Activated field to false, which isn't strictly necessary because the
	// Activated field will have the zero-value of false by default. But setting this
	// explicitly helps to make our intentions clear to anyone reading the code.
	user := &model.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}
	// Use the Password.Set() method to generate and store the hashed and plaintext
	// passwords.
	err = user.Password.Set(input.Password)
	if err != nil {
		api.serverErrorResponse(w, r, err)
		return
	}
	v := validator.New()
	// Validate the user struct and return the error messages to the client if any of
	// the checks fail.
	if model.ValidateUser(v, user); !v.Valid() {
		api.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Insert the user data into the database.
	err = api.UserModel.Insert(user)
	if err != nil {
		switch {
		// If we get a ErrDuplicateEmail error, use the v.AddError() method to manually
		// add a message to the validator instance, and then call our
		// failedValidationResponse() helper.
		case errors.Is(err, model.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			api.failedValidationResponse(w, r, v.Errors)
		default:
			api.serverErrorResponse(w, r, err)
		}
		return
	}
	//err = api.PermissionModel.AddForUser(user.ID, "shop:read")
	//if err != nil {
	//	api.serverErrorResponse(w, r, err)
	//	return
	//}

	err = api.PermissionModel.AddForUser(user.ID, "shop:read")
	if err != nil {
		api.serverErrorResponse(w, r, err)
		return
	}

	// After the user record has been created in the database, generate a new activation
	// token for the user.
	token, err := api.TokenModel.New(user.ID, 3*24*time.Hour, model.ScopeActivation)
	if err != nil {
		api.serverErrorResponse(w, r, err)
		return
	}

	var res struct {
		Token *string     `json:"token"`
		User  *model.User `json:"user"`
	}

	res.Token = &token.Plaintext
	res.User = user

	api.writeJSON(w, http.StatusCreated, envelope{"user": res}, nil)
}

func (api *API) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the plaintext activation token from the request body
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

	// Retrieve the details of the user associated with the token using the GetForToken() method.
	// If no matching record is found, then we let the client know that the token they provided
	// is not valid.
	user, err := api.UserModel.GetForToken(model.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			api.failedValidationResponse(w, r, v.Errors)
		default:
			api.serverErrorResponse(w, r, err)
		}
		return
	}

	// Update the user's activation status.
	user.Activated = true

	// Save the updated user record in our database, checking for any edit conflicts in the same
	// way that we did for our move records.
	err = api.UserModel.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrEditConflict):
			api.editConflictResponse(w, r)
		default:
			api.serverErrorResponse(w, r, err)
		}
		return
	}

	// If everything went successfully above, then delete all activation tokens for the user.
	err = api.TokenModel.DeleteAllForUser(model.ScopeActivation, user.ID)
	if err != nil {
		api.serverErrorResponse(w, r, err)
		return
	}

	api.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
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
	token, err = api.TokenModel.New(user.ID, 24*time.Hour, model.ScopeAuthentication)
	if err != nil {
		api.serverErrorResponse(w, r, err)
		return
	}

	// Encode the token to JSON and send it in the response along with a 201 Created status code.
	err = api.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
		api.serverErrorResponse(w, r, err)
	}
}

type contextKey string

// userContextKey is used as a key for getting and setting user information in the request
// context.
const userContextKey = contextKey("user")

// contextSetUser returns a new copy of the request with the provided User struct added to the
// context.
func (api *API) contextSetUser(r *http.Request, user *model.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// contextGetUser retrieves the User struct from the request context. The only time that
// this helper should be used is when we logically expect there to be a User struct value
// in the context, and if it doesn't exist it will firmly be an 'unexpected' error, upon we panic.
func (api *API) contextGetUser(r *http.Request) *model.User {
	user, ok := r.Context().Value(userContextKey).(*model.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
