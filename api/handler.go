package api

import (
	"Golang_Project/pkg/model"

	//"errors"
	"fmt"
	"log"
	"net/http"
	//"time"
)

type Response struct {
	Shops       []model.Shop        `json:"shops"`
	Products    []model.Product     `json:"products"`
	Users       []model.User        `json:"users"`
	Tokens      []model.Token       `json:"tokens"`
	Permissions []model.Permissions `json:"permissions"`
	Carts       []model.Cart        `json:"carts"`
}

type API struct {
	ShopModel       *model.ShopModel
	ProductModel    *model.ProductModel
	UserModel       *model.UserModel
	TokenModel      *model.TokenModel
	PermissionModel *model.PermissionModel
	CartModel       *model.CartModel
}

func NewAPI(shopModel *model.ShopModel, productModel *model.ProductModel, userModel *model.UserModel, tokenModel *model.TokenModel, permissionModel *model.PermissionModel, cartModel *model.CartModel) *API {
	return &API{ShopModel: shopModel, ProductModel: productModel, UserModel: userModel, TokenModel: tokenModel, PermissionModel: permissionModel, CartModel: cartModel}
}

func (api *API) HealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("welcome")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello there")
}
