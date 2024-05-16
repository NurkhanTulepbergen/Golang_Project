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
	Orders      []model.Order       `json:"orders"`
	History     []model.History     `json:"history"`
}

type API struct {
	ShopModel       *model.ShopModel
	ProductModel    *model.ProductModel
	UserModel       *model.UserModel
	TokenModel      *model.TokenModel
	PermissionModel *model.PermissionModel
	CartModel       *model.CartModel
	OrderModel      *model.OrderModel
	HistoryModel    *model.HistoryModel
}

func NewAPI(shopModel *model.ShopModel, productModel *model.ProductModel, userModel *model.UserModel, tokenModel *model.TokenModel, permissionModel *model.PermissionModel, cartModel *model.CartModel, orderModel *model.OrderModel, historyModel *model.HistoryModel) *API {
	return &API{ShopModel: shopModel, ProductModel: productModel, UserModel: userModel, TokenModel: tokenModel, PermissionModel: permissionModel, CartModel: cartModel, OrderModel: orderModel, HistoryModel: historyModel}
}

func (api *API) HealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("welcome")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello there")
}
