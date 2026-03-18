package main

import (
	"github.com/klyakssa/go-diplom.git/internal/app"
	"github.com/klyakssa/go-diplom.git/internal/config"
	"github.com/shopspring/decimal"
)

func main() {
	config := config.InitConfiguration()
	decimal.MarshalJSONWithoutQuotes = true

	app.Run(config)
}
