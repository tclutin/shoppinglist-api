package main

import (
	"context"
	"github.com/tclutin/shoppinglist-api/internal/app"
)

func main() {
	app.New().Run(context.Background())
}
