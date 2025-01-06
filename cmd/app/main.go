package main

import (
	"context"
	_ "github.com/tclutin/shoppinglist-api/docs"
	"github.com/tclutin/shoppinglist-api/internal/app"
)

//	@title			ShoppingList API
//	@version		1.0
//	@description	for pet project for my friend

//	@host		localhost:8081
//	@BasePath	/api/

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				Use "Bearer <token>" to authenticate
func main() {
	app.New().Run(context.Background())
}
